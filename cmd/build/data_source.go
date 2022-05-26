package build

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/plentico/plenti/common"
	"github.com/plentico/plenti/readers"
	"github.com/spf13/afero"
)

// Local is set to true when using dev webserver, otherwise bools default to false.
var Local bool

// Doreload and other flags should probably be part of a config accessible across build.
// It gets set using server flags.
var Doreload bool
var (
	// Setup regex to find field name.
	reField = regexp.MustCompile(`:field\((.*?)\)`)
	// Setup regex to find pagination and a leading forward slash.
	rePaginate = regexp.MustCompile(`:paginate\((.*?)\)`)

	// Create regex for allowed characters when slugifying path.
	reSlugify = regexp.MustCompile("[^a-z0-9/]+")
	// Remove newlines.
	reN = regexp.MustCompile(`\r?\n`)

	// Remove tabs.
	reT = regexp.MustCompile(`\t`)

	// Reduce extra whitespace to a single space.
	reS = regexp.MustCompile(`\s+`)
)

// Holds info related to a particular content node.
type content struct {
	contentType      string
	contentPath      string
	contentDest      string
	contentDetails   string
	contentFilepath  string
	contentFilename  string
	contentFields    string
	contentPagerDest string
	contentPagerPath string
	contentPagerNums []string
}

// Holds sitewide environment variables.
type env struct {
	Local   bool   `json:"local"`
	Baseurl string `json:"baseUrl"`
	Cms     cms    `json:"cms"`
}
type cms struct {
	Repo        string              `json:"repo"`
	RedirectUrl string              `json:"redirectUrl"`
	AppId       string              `json:"appId"`
	Branch      string              `json:"branch"`
	ImageSizes  map[string][][2]int `json:"imageSizes"`
}

// DataSource builds json list from "content/" directory.
func DataSource(buildPath string, siteConfig readers.SiteConfig) error {

	defer Benchmark(time.Now(), "Creating data_source")

	Log("\nGathering data source from 'content/' folder")

	// Set some defaults
	contentJSPath := buildPath + "/spa/ejected/content.js"
	envPath := buildPath + "/spa/ejected/env.js"
	env := env{
		Local:   Local,
		Baseurl: siteConfig.BaseURL,
		Cms: cms{
			Repo:        siteConfig.CMS.Repo,
			RedirectUrl: siteConfig.CMS.RedirectUrl,
			AppId:       siteConfig.CMS.AppId,
			Branch:      siteConfig.CMS.Branch,
			ImageSizes:  siteConfig.CMS.ImageSizes,
		},
	}

	// Create env magic prop.
	envJson, err := json.Marshal(env)
	if err != nil {
		fmt.Printf("Unable to marshal env to JSON")
		return err
	}
	envStr := "export let env = " + string(envJson) + ";"

	// no dirs needed for mem
	if common.UseMemFS {
		common.Set(contentJSPath, "", &common.FData{B: []byte(`const allContent = [`)})
		common.Set(envPath, "", &common.FData{B: []byte(envStr)})
	} else {
		if err := os.MkdirAll(buildPath+"/spa/ejected", os.ModePerm); err != nil {
			return err
		}
		// Start the new content.js file.
		err := ioutil.WriteFile(contentJSPath, []byte(`const allContent = [`), 0755)
		if err != nil {
			fmt.Printf("Unable to write content.js file: %v", err)
			return err
		}
		// Create the env.js file.
		err = ioutil.WriteFile(envPath, []byte(envStr), 0755)
		if err != nil {
			fmt.Printf("Unable to write env.js file: %v", err)
			return err
		}
	}

	// Set up counter for logging output.
	contentFileCounter := 0
	// Start the string that will be used for allContent object.
	allContentStr := "["
	// Store each content file in array we can iterate over for creating static html.
	allContent := []content{}

	// Go through all sub directories in "content/" folder.
	if ThemeFs != nil {
		if err := afero.Walk(ThemeFs, "content", func(path string, info os.FileInfo, err error) error {
			contentFileCounter, allContentStr, allContent, err = getContent(path, info, err, siteConfig, buildPath, contentJSPath, allContentStr, allContent, contentFileCounter)
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			return fmt.Errorf("Could not get layout %w from virtual theme\n", err)
		}
	} else {
		if err := filepath.Walk("content", func(path string, info os.FileInfo, err error) error {
			contentFileCounter, allContentStr, allContent, err = getContent(path, info, err, siteConfig, buildPath, contentJSPath, allContentStr, allContent, contentFileCounter)
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			return fmt.Errorf("Could not get layout %w from project\n", err)
		}
	}

	// End the string that will be used in allContent object.
	allContentStr = strings.TrimSuffix(allContentStr, ",") + "]"

	for _, currentContent := range allContent {

		if err := createProps(currentContent, allContentStr, string(envJson)); err != nil {
			return err
		}

		if err := createHTML(currentContent); err != nil {
			return err
		}

		allPaginatedContent, err := paginate(currentContent, contentJSPath)
		if err != nil {
			return err
		}
		for _, paginatedContent := range allPaginatedContent {
			if err = createProps(paginatedContent, allContentStr, string(envJson)); err != nil {
				return err
			}

			if err = createHTML(paginatedContent); err != nil {
				return err
			}

		}

	}

	Log("Number of content files used: " + fmt.Sprint(contentFileCounter))
	// Complete the content.js file.
	if err := writeContentJS(contentJSPath, "];\n\nexport default allContent;"); err != nil {
		return err
	}
	if common.UseMemFS {
		// hash all content
		common.Get(contentJSPath).Hash = common.CRC32Hasher(common.Get(contentJSPath).B)
	}
	return nil

}

func getContent(path string, info os.FileInfo, err error, siteConfig readers.SiteConfig,
	buildPath string, contentJSPath string, allContentStr string, allContent []content,
	contentFileCounter int) (int, string, []content, error) {
	if err != nil {
		return contentFileCounter, allContentStr, allContent, fmt.Errorf("can't stat %s: %w", path, err)
	}
	if !info.IsDir() {
		filePath := path
		// Get individual path arguments.
		parts := strings.Split(path, "/")
		contentType := parts[1]
		fileName := parts[len(parts)-1]

		// Don't add _blueprint.json or other special named files starting with underscores.
		if fileName[:1] != "_" && fileName[:1] != "." {

			// Get the contents of the file.
			fileContentBytes, err := getVirtualFileIfThemeBuild(path)
			if err != nil {
				return contentFileCounter, allContentStr, allContent, fmt.Errorf("file: %s %w%s\n", path, err, common.Caller())
			}
			fileContentStr := string(fileContentBytes)

			// Remove the "content/" folder from path.
			path = strings.TrimPrefix(path, "content/")

			// Check for index file at any level.
			if fileName == "index.json" {
				// Remove entire filename from path.
				path = strings.TrimSuffix(path, fileName)
			} else {
				// Remove file extension only from path for files other than index.json.
				path = strings.TrimSuffix(path, filepath.Ext(path))
			}

			// Remove the extension (if it exists) from single types since the filename = the type name.
			contentType = strings.TrimSuffix(contentType, filepath.Ext(contentType))

			// Get field key/values from content source.
			typeFields := readers.GetTypeFields(fileContentBytes)
			// Setup regex to find field name.
			reField := regexp.MustCompile(`:fields\((.*?)\)`)
			// Check for path overrides from plenti.json config file.
			for configContentType, slug := range siteConfig.Routes {
				if configContentType == contentType {
					// Replace :filename.
					slug = strings.Replace(slug, ":filename", strings.TrimSuffix(fileName, filepath.Ext(fileName)), -1)

					// Replace :fields().
					fieldReplacements := reField.FindAllStringSubmatch(slug, -1)
					// Loop through all :fields() replacements found in config file.
					for _, replacement := range fieldReplacements {
						// Loop through all top level keys found in content source file.
						for field, fieldValue := range typeFields.Fields {
							// Check if field name in the replacement pattern is found in data source.
							if replacement[1] == field {
								// Use the field value in the path.
								slug = strings.ReplaceAll(slug, replacement[0], fieldValue)
							}
						}
					}
					path = slug
				}
			}

			// Initialize vars for path with replacement patterns still intact.
			var pagerPath string
			var pagerDestPath string
			// If there is a /:paginate() replacement found.
			if rePaginate.MatchString(path) {
				// Save path before slugifying to preserve pagination.
				pagerPath = path
				// Get Destination path before slugifying to preserve pagination.
				pagerDestPath = buildPath + "/" + path + "/index.html"
				// Remove /:pagination()
				path = rePaginate.ReplaceAllString(path, "")
			}

			// Slugify output using reSlugify regex defined above.
			path = strings.Trim(reSlugify.ReplaceAllString(strings.ToLower(path), "-"), "-")

			// Check for any missing/blank paths.
			if path == "" {
				// Add the forward slash back for the index page.
				path = "/"
			}
			// Let the user know if path is blank.
			if len(path) < 1 {
				fmt.Println("Content path can't be blank, check your route overrides in plenti.json.")
			}

			// Remove any repeating forward slashes.
			path = filepath.Clean(path)

			// Remove trailing slash, unless it's the homepage.
			if path != "/" && path[len(path)-1:] == "/" {
				path = strings.TrimSuffix(path, "/")
			}

			destPath := buildPath + "/" + path + "/index.html"

			contentDetailsStr := "{\n" +
				"\"pager\": 1,\n" +
				"\"type\": \"" + contentType + "\",\n" +
				"\"path\": \"" + path + "\",\n" +
				"\"filepath\": \"" + filePath + "\",\n" +
				"\"filename\": \"" + fileName + "\",\n" +
				"\"fields\": " + fileContentStr + "\n}"

			// Write to the content.js client data source file.
			if err = writeContentJS(contentJSPath, contentDetailsStr+","); err != nil {
				return contentFileCounter, allContentStr, allContent, fmt.Errorf("file: %s %w%s\n", contentJSPath, err, common.Caller())
			}

			// Remove newlines, tabs, and extra space.
			encodedContentDetails := encodeString(contentDetailsStr)
			// Add info for being referenced in allContent object.
			allContentStr = allContentStr + encodedContentDetails + ","

			content := content{
				contentType:      contentType,
				contentPath:      path,
				contentDest:      destPath,
				contentDetails:   encodedContentDetails,
				contentFilepath:  filePath,
				contentFilename:  fileName,
				contentFields:    encodeString(fileContentStr),
				contentPagerDest: pagerDestPath,
				contentPagerPath: pagerPath,
			}
			allContent = append(allContent, content)

			// Increment counter for logging purposes.
			contentFileCounter++

		}
	}
	return contentFileCounter, allContentStr, allContent, nil
}

func createProps(currentContent content, allContentStr string, env string) error {
	componentSignature := "layouts_content_" + currentContent.contentType + "_svelte"
	_, err := SSRctx.RunScript("var props = {content: "+currentContent.contentDetails+
		", layout: "+componentSignature+
		", allContent: "+allContentStr+
		", env: "+env+
		"};", "create_ssr")
	if err != nil {
		return fmt.Errorf("Could not create props: %w%s\n", err, common.Caller())
	}
	// Render the HTML with props needed for the current content.
	_, err = SSRctx.RunScript("var { html, css: staticCss} = layouts_global_html_svelte.render(props);", "create_ssr")
	if err != nil {
		return fmt.Errorf("Can't render htmlComponent: %w%s\n", err, common.Caller())

	}
	return nil
}

func createHTML(currentContent content) error {
	// Get the rendered HTML from v8go.
	renderedHTML, err := SSRctx.RunScript("html;", "create_ssr")
	if err != nil {
		return fmt.Errorf("V8go could not execute js default: %w%s\n", err, common.Caller())

	}
	// Get the string value of the static HTML.
	renderedHTMLStr := renderedHTML.String()
	// Convert the string to byte array that can be written to file system.
	htmlBytes := []byte(renderedHTMLStr)
	// Inject <!DOCTYPE html>
	htmlBytes = bytes.Replace(htmlBytes, []byte("<html"), []byte("<!DOCTYPE html><html"), 1)
	if Doreload {
		// Inject live-reload script (stored in ejected core).
		htmlBytes = bytes.Replace(htmlBytes, []byte("</body>"), []byte("<script type='text/javascript' src='/spa/ejected/live-reload.js'></script></body>"), 1)
	}
	if common.UseMemFS {
		common.Set(currentContent.contentDest, "", &common.FData{B: htmlBytes})
		return nil
	}
	// Create any folders need to write file.
	if err := os.MkdirAll(strings.TrimSuffix(currentContent.contentDest, "/index.html"), os.ModePerm); err != nil {
		return fmt.Errorf("couldn't create dirs in createHTML: %w%s\n", err, common.Caller())
	}
	// Write static HTML to the filesystem.
	err = ioutil.WriteFile(currentContent.contentDest, htmlBytes, 0755)
	if err != nil {
		return fmt.Errorf("unable to write SSR file: %w%s\n", err, common.Caller())
	}
	return nil
}

func paginate(currentContent content, contentJSPath string) ([]content, error) {
	paginatedContent, _ := getPagination()
	var err error
	allNewContent := []content{}
	// Loop through all :paginate() replacements found in config file.
	for _, pager := range paginatedContent {
		// Check if the config file specifies pagination for this Type.
		if len(pager.paginationVars) > 0 && pager.contentType == currentContent.contentType {
			// Increment the pager.
			allNewContent, err = incrementPager(pager.paginationVars, currentContent, contentJSPath, allNewContent)
			if err != nil {
				return nil, err
			}
		}
	}
	return allNewContent, err
}

func incrementPager(paginationVars []string, currentContent content, contentJSPath string, allNewContent []content) ([]content, error) {
	// Pop first item from the list.
	paginationVar, paginationVars := paginationVars[0], paginationVars[1:]
	// Copy the current content so we can increment the pager.
	newContent := currentContent
	// Get the number of pages for the pager.
	totalPagesInt, err := getTotalPages(paginationVar)
	if err != nil {
		return nil, err
	}
	// Loop through total number of pages for current pager.
	for i := 1; i <= totalPagesInt; i++ {
		// Convert page number to a string that can be used in a path.
		currentPageNumber := strconv.Itoa(i)
		// Update the path by replacing the current :paginate() pattern.
		newContent.contentPath = strings.Replace(currentContent.contentPagerPath, ":paginate("+paginationVar+")", currentPageNumber, 1)
		newContent.contentDest = strings.Replace(currentContent.contentPagerDest, ":paginate("+paginationVar+")", currentPageNumber, 1)
		// Now we need to update the pager path to replace the pattern with a number so it's ready if we call incrementPager() again for a second pager.
		newContent.contentPagerPath = newContent.contentPath
		newContent.contentPagerDest = newContent.contentDest

		// Collect each pager value.
		newContent.contentPagerNums = append(newContent.contentPagerNums, currentPageNumber)

		// Check if there are more pagers for the route override.
		if len(paginationVars) > 0 {
			// Recursively call func to increment second pager.
			// todo: a better approach
			allNewContentTmp, err := incrementPager(paginationVars, newContent, contentJSPath, allNewContent)
			if err != nil {
				return nil, err
			}
			allNewContent = allNewContentTmp
			// Remove first item in the array to move to the next number in the first pager.
			newContent.contentPagerNums = newContent.contentPagerNums[1:]
			// Continue because you don't want to complete the loop with a partially updated path (we have more pagers!).
			continue
		}

		// Set the content.pager value if only one pager is used.
		pageNums := newContent.contentPagerNums[0]
		// Check if there are multiple pagers in the router override.
		if len(newContent.contentPagerNums) > 1 {
			// Make the content.pager value an array with the current pager values.
			pageNums = "[" + strings.Join(newContent.contentPagerNums, ", ") + "]"
		}

		// Add current page number to the content source so it can be pulled in as the current page.
		newContent.contentDetails = "{\n" +
			"\"pager\": " + pageNums + ",\n" +
			"\"type\": \"" + newContent.contentType + "\",\n" +
			"\"path\": \"" + newContent.contentPath + "\",\n" +
			"\"filepath\": \"" + newContent.contentFilepath + "\",\n" +
			"\"filename\": \"" + newContent.contentFilename + "\",\n" +
			"\"fields\": " + newContent.contentFields + "\n}"

		// Add paginated entries to content.js file.
		if err := writeContentJS(contentJSPath, newContent.contentDetails+","); err != nil {
			return nil, err
		}
		// Add to array of content for creating paginated static HTML fallbacks.
		allNewContent = append(allNewContent, newContent)

		// Remove last number from array to get next page in current pager.
		newContent.contentPagerNums = newContent.contentPagerNums[:len(newContent.contentPagerNums)-1]
	}
	return allNewContent, nil
}

func getTotalPages(paginationVar string) (int, error) {
	totalPages, err := SSRctx.RunScript("plenti_global_pager_"+paginationVar, "create_ssr")
	if err != nil {
		return 0, fmt.Errorf("Could not get value of '%v' used in pager: %w%s\n", paginationVar, err, common.Caller())
	}
	// Convert string total page value to integer.
	totalPagesInt, err := strconv.Atoi(totalPages.String())
	if err != nil {
		return 0, fmt.Errorf("Can't convert pager value '%s' to an integer: %w%s\n", totalPages.String(), err, common.Caller())
	}
	return totalPagesInt, nil
}

func writeContentJS(contentJSPath string, contentDetailsStr string) error {
	if common.UseMemFS {
		allB := append(common.Get(contentJSPath).B, []byte(contentDetailsStr)...)
		// ok to append as it gets created each build and has
		common.Set(contentJSPath, "", &common.FData{B: allB})
		return nil
	}
	// Create new content.js file if it doesn't already exist, or add to it if it does.
	contentJSFile, err := os.OpenFile(contentJSPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open content.js for writing: %w%s\n", err, common.Caller())
	}
	// Write to the file with info from current file in "/content" folder.
	defer contentJSFile.Close()
	if _, err := contentJSFile.WriteString(contentDetailsStr); err != nil {

		return fmt.Errorf("could not write to file %s: %w%s\n", contentJSPath, err, common.Caller())

	}
	return nil
}

func encodeString(encodedStr string) string {
	// Remove newlines.
	encodedStr = reN.ReplaceAllString(encodedStr, " ")
	// Remove tabs.

	encodedStr = reT.ReplaceAllString(encodedStr, " ")
	// Reduce extra whitespace to a single space.

	encodedStr = reS.ReplaceAllString(encodedStr, " ")
	return encodedStr
}
