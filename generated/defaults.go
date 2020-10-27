package generated

// Do not edit, this file is automatically generated.

// Defaults: scaffolding used in 'build' command
var Defaults = map[string][]byte{
	"/.gitignore": []byte(`public
node_modules`),
	"/assets/logo.svg": []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg width="40" height="40" version="1.1" viewBox="0 0 10.583 10.583" xmlns="http://www.w3.org/2000/svg" xmlns:cc="http://creativecommons.org/ns#" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
 <metadata>
  <rdf:RDF>
   <cc:Work rdf:about="">
    <dc:format>image/svg+xml</dc:format>
    <dc:type rdf:resource="http://purl.org/dc/dcmitype/StillImage"/>
    <dc:title/>
   </cc:Work>
  </rdf:RDF>
 </metadata>
 <g transform="translate(63.525 -69.957)" fill-rule="evenodd">
  <path d="m-54.57 74.523c-1.7919 1.907-5.3906 2.3551-4.3001 4.9881 0.36199 0.74862-0.27419 1.0155-0.50107 0.61951 0 0-3.4086-2.8786-1.4504-4.6156 0.98084-1.0161 6.9801-5.3003 6.3041-1.3029z" fill="#1c7fc7"/>
  <path d="m-57.629 73.238s-1.5075 0.83503-2.7039 1.8531c0 0 3.0859-1.3869 5.7631-0.56799 0 0 0.39145-1.6984-3.0593-1.2851z" fill="#004a77" fill-opacity=".87434"/>
  <path d="m-62.252 72.183s0.38056-1.1878 2.1798-1.9508l0.96418 0.28985s6.7226-1.3794 4.5374 4.0016c0 0 0.5771-2.071-4.5695-1.033l-0.70199 0.43123c-1.906-0.26683-2.4099-1.7389-2.4099-1.7389" fill="#22a6ed"/>
 </g>
</svg>
`),
	"/content/blog/_blueprint.json": []byte(`{
    "title": "text",
    "body": ["text"],
    "author": "text",
    "date": "date"
}`),
	"/content/blog/perry.json": []byte(`{
    "title": "Customize your Planarian",
    "body": [
        "Meet our mascot - Perry the Planarian!",
        "If you want to customize Perry's style, you can do so <a href='https://perry.plenti.co' target='blank' rel='noopener noreferrer'>here</a> (Coming soon).",
        "Check out this <a href='https://youtu.be/Gr3KTOnsWEM' target='blank' rel='noopener noreferrer'>video</a> to see how the app was made!"
    ],
    "author": "Jim Fisk",
    "date": "10/1/2020"
}`),
	"/content/blog/pletiform.json": []byte(`{
    "title": "Build sites with good form",
    "body": [
        "Need an easy webform solution?",
        "Try adding a <a href='https://plentiform.com' target='blank' rel='noopener noreferrer'>plentiform</a>! (Coming soon)"
    ],
    "author": "Jim Fisk",
    "date": "1/26/2020"
}`),
	"/content/blog/stores.json": []byte(`{
    "title": "Svelte writable stores example",
    "body": [
        "Here's a basic counter implementation using Svelte Writable Stores"
    ],
    "store": true,
    "author": "Jim Fisk",
    "date": "8/25/2020"
}`),
	"/content/index.json": []byte(`{
	"title": "Welcome to Plenti!",
	"intro": [
		"Take a look around to see how things work.",
		"The bottom of each page will tell you where to find the corresponding template in your project.",
		"If you get stuck, check out our <a href='https://plenti.co/docs' target='blank' rel='noopener noreferrer'>docs</a>. If you need extra help, <a href='https://github.com/plentico/plenti/issues/new' target='blank' rel='noopener noreferrer'>let us know</a>! Enjoy :)"
	],
	"components": [
		{
			"component": "template",
			"fields": {"type": "index"}
		}
	]
}`),
	"/content/pages/_blueprint.json": []byte(`{
	"title": "text",
	"description": ["text"],
	"author": "text"
}`),
	"/content/pages/about.json": []byte(`{
	"title": "About Plenti",
	"description": [
		"Plenti is a minimalist <a href='https://jamstack.org/' target='blank' rel='noopener noreferrer'>JAMstack</a> framework that's flexible and easy to use.",
		"We've cut out as many dependencies as possible so you can focus on being productive instead of wrestling with tools.",
		"The <a href='https://svelte.dev/' target='blank' rel='noopener noreferrer'>Svelte</a> frontend gives users get a snappy experience, and the <a href='https://golang.org/' target='blank' rel='noopener noreferrer'>Go</a> backend builds fast so you can get more done."
	],
	"author": "Jim Fisk"
}`),
	"/content/pages/contact.json": []byte(`{
	"title": "Contact Us",
	"description": [
		"Plenti is 100% free and open source!",
		"You can fork it for your own purposes, or help us out by reporting bugs / contributing code on <a href='https://github.com/plentico/plenti' target='blank' rel='noopener noreferrer'>Our GitHub</a>.",
		"Give us <a href='https://twitter.com/plentico' target='blank' rel='noopener noreferrer'>a tweet</a> if you like what you see!"
	],
	"author": "Jim Fisk"
}`),
	"/layout/components/decrementer.svelte": []byte(`<script>
  import { count } from '../scripts/stores.svelte';

  function decrement() {
    count.update(n => n - 1);
  }
</script>

<button on:click={decrement}>
  -
</button>`),
	"/layout/components/grid.svelte": []byte(`<script>
  import { sortByDate } from '../scripts/sort_by_date.svelte';
  export let items, filter;
</script>

<div class="grid">
  {#each sortByDate(items) as item}
		{#if item.type == filter}
      <a class="grid-item" href="{item.path}">{item.fields.title}</a>
		{/if}
  {/each}
</div>

<style>
  .grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    grid-column-gap: 10px;
    grid-row-gap: 10px;
  }
  .grid-item {
    display: flex;
    flex-grow: 1;
    height: 200px;
    align-items: center;
    justify-content: center;
    background: var(--primary);
    font-weight: bold;
    border-radius: 5px;
    color: white;
  }
  a::before {
    content: none;
  }
</style>
`),
	"/layout/components/incrementer.svelte": []byte(`<script>
  import { count } from '../scripts/stores.svelte';

  function increment() {
    count.update(n => n + 1);
  }
</script>

<button on:click={increment}>
  +
</button>`),
	"/layout/components/template.svelte": []byte(`<script>
  export let type;

  let path;
  let copyText = "Copy";
  const copy = async () => {
    if (!navigator.clipboard) {
      return
    }
    try {
      copyText = "Copied";
      await navigator.clipboard.writeText(path.innerHTML);
      setTimeout(() => copyText = "Copy", 500);
    } catch (err) {
      console.error('Failed to copy!', err)
    }
  }
</script>

<div class="template">
  <span>Template:</span>
  <pre>
    <code bind:this={path} class="{copyText}">layout/content/{type}.svelte</code>
    <button on:click={copy}>{copyText}</button>
  </pre>
</div>

<style>
  .template {
    display: flex;
    align-items: center;
  }
  pre {
    display: flex;
    padding-left: 5px;
  }
  code {
      background-color: var(--base);
      padding: 5px 10px;
  }
  code.copied {
      background-color: var(--accent);
  }
  button {
    border: 1px solid rgba(0,0,0,.1);
    background: white;
    padding: 4px;
    border-top-right-radius: 5px;
    border-bottom-right-radius: 5px;
    cursor: pointer;
  }
</style>`),
	"/layout/content/404.svelte": []byte(`<h1>Oops... 404 not found</h1>
<a href="/">Go home?</a>`),
	"/layout/content/blog.svelte": []byte(`<script>
	export let title, body, author, date, store;
  import Uses from "../components/template.svelte";

  // Svelte store example:
  import { count } from '../scripts/stores.svelte';
  import Incrementer from '../components/incrementer.svelte';
  import Decrementer from '../components/decrementer.svelte';
  let count_value;
  const unsubscribe = count.subscribe(value => {
    count_value = value;
  });
</script>

<h1>{title}</h1>

<p><em>{#if author}Written by {author}{/if}{#if date}&nbsp;on {date}{/if}</em></p>

<div>
  {#each body as paragraph}
    <p>{@html paragraph}</p>
  {/each}
</div>

{#if store}
  <h3>The count is {count_value}</h3>
  <Incrementer/>
  <Decrementer/>  
{/if}

<Uses type="blog" />

<p><a href="/">Back home</a></p>
`),
	"/layout/content/index.svelte": []byte(`<script>
	export let title, intro, components, allContent;
	import Grid from '../components/grid.svelte';
	import { loadComponent } from '../scripts/load_component.svelte';
</script>

<h1>{title}</h1>

<section id="intro">
	{#each intro as paragraph}
		<p>{@html paragraph}</p>
	{/each}
</section>

<div>
	<h3>Recent blog posts:</h3>
	<Grid items={allContent} filter="blog" />
	<br />
</div>

{#if components}
	{#each components as { component, fields }}
		{#await loadComponent(component)}
			loading component...
		{:then compClass}
			<svelte:component this="{compClass}" {...fields} />
		{:catch error}
			{console.log(error.message)}
		{/await}
	{/each}
{/if}`),
	"/layout/content/pages.svelte": []byte(`<script>
  export let title, description;
  import Uses from "../components/template.svelte";
</script>

<h1>{title}</h1>

<div>
  {#each description as paragraph}
    <p>{@html paragraph}</p>
  {/each}
</div>

<Uses type="pages" />

<p><a href="/">Back home</a></p>`),
	"/layout/global/footer.svelte": []byte(`<script>
  export let allContent;
  import { makeTitle } from '../scripts/make_title.svelte';
</script>

<footer>
  <div class="container">
    <span>All content:</span>
    {#each allContent as content}
      <a href="{content.path}">{makeTitle(content.filename)}</a>
    {/each}
  </div>
</footer>

<style>
  footer {
    min-height: 200px;
    display: flex;
    align-items: center;
    background-color: var(--base-dark);
    margin-top: 100px;
  }
  span {
    color: var(--primary);
    font-weight: bold;
  }
  a {
    color: white;
    text-decoration: none;
    margin-left: 10px;
  }
</style>
`),
	"/layout/global/head.svelte": []byte(`<script>
  export let title;
</script>

<head>
  <meta charset='utf-8'>
  <meta name='viewport' content='width=device-width,initial-scale=1'>

  <title>{ title }</title>

  <link href="https://fonts.googleapis.com/css2?family=Rubik:ital,wght@0,300;0,700;1,300&display=swap" rel="stylesheet">
  <link rel="icon" type="image/svg+xml" href="/assets/logo.svg">
  <link rel='stylesheet' href='/spa/bundle.css'>
</head>
`),
	"/layout/global/html.svelte": []byte(`<script>
  import Head from './head.svelte';
  import Nav from './nav.svelte';
  import Footer from './footer.svelte';
  import { makeTitle } from '../scripts/make_title.svelte';

  export let route, content, allContent;
</script>

<html lang="en">
<Head title={makeTitle(content.filename)} />
<body>
  <Nav />
  <main>
    <div class="container">
      <svelte:component this={route} {...content.fields} {allContent} />
      <br />
    </div>
  </main>
  <Footer {allContent} />
</body>
</html>

<style>
  body {
    font-family: 'Rubik', sans-serif;
    display: flex;
    flex-direction: column;
    margin: 0;
  }
  main {
    flex-grow: 1;
  }
  :global(.container) {
    max-width: 1024px;
    margin: 0 auto;
    flex-grow: 1;
    padding: 0 20px;
  }
  :global(:root) {
    --primary: rgb(34, 166, 237);
    --primary-dark: rgb(16, 92, 133);
    --accent: rgb(254, 211, 48);
    --base: rgb(245, 245, 245);
    --base-dark: rgb(17, 17, 17);
  }
  :global(main a) {
    position: relative;
    text-decoration: none;
    color: var(--base-dark);
    padding-bottom: 5px;
  }
  :global(main a:before) {
    content: "";
    width: 100%;
    height: 100%;
    background-image: linear-gradient(to top, var(--accent) 25%, rgba(0, 0, 0, 0) 40%);  
    position: absolute;
    left: 0;
    bottom: 2px;
    z-index: -1;   
    will-change: width;
    transform: rotate(-2deg);
    transform-origin: left bottom;
    transition: width .1s ease-out;
  }
  :global(main a:hover:before) {
    width: 0;
    transition-duration: .15s;
  }
</style>
`),
	"/layout/global/nav.svelte": []byte(`<nav>
  <div class="container">
    <span id="brand"><a href="/"><img alt="planarian" src="/assets/logo.svg" />Home</a></span>
    <a href="/about">About</a>&nbsp;
    <a href="/contact">Contact</a>
  </div>
</nav>

<style>
  nav, .container, #brand, #brand a {
    display: flex;
  }
  nav {
    min-height: 60px;
    box-shadow: 0px 2px 3px var(--base);
  }
  #brand {
    flex: 1;
  }
  a {
    align-self: center;
    align-items: center;
    color: var(--base-dark);
    text-decoration: none;
  }
  img {
    margin-right: 10px;
  }
</style>
`),
	"/layout/scripts/load_component.svelte": []byte(`<script context="module">
  export const loadComponent = component => {
    let compClassPromise = import("../components/" + component + ".svelte").then(res => res.default);
    // Fix "Unhandled promise rejection" error.
    // See: https://github.com/sveltejs/sapper/issues/487#issuecomment-529145749
    compClassPromise.catch(err => null)
    return compClassPromise;
  }
</script>
`),
	"/layout/scripts/make_title.svelte": []byte(`<script context="module">
  export const makeTitle = filename => {
  if (filename == 'index.json') {
    return 'Home';
  } else if (filename) {
    // Remove file extension.
    filename = filename.split('.').slice(0, -1).join('.');
    // Convert underscores and hyphens to spaces.
    filename = filename.replace(/_|-/g,' ');
    // Capitalize first letter of each word.
    filename = filename.split(' ').map((s) => s.charAt(0).toUpperCase() + s.substring(1)).join(' ');
  } 
  return filename;
  }
</script>
`),
	"/layout/scripts/sort_by_date.svelte": []byte(`<script context="module">
  export const sortByDate = (items, order) => {
    items.sort((a, b) => { 
      // Must have a field specifically named "date" to work.
      // Feel free to extend to other custom named date fields.
      if (a.fields.hasOwnProperty("date") && b.fields.hasOwnProperty("date")) {
        let aDate = new Date(a.fields.date);
        let bDate = new Date(b.fields.date);
        if (order == "oldest") {
            return aDate - bDate;
        }
        return bDate - aDate;
      }
    });
    return items;
  };
</script>`),
	"/layout/scripts/stores.svelte": []byte(`<script context="module">
  import { writable } from 'svelte/store';
  export const count = writable(0);
</script>`),
	"/package.json": []byte(`{
  "name": "my-plenti-app",
  "version": "1.0.0",
  "type": "module",
  "private": true,
  "dependencies": {
    "navaid": "^1.1.1",
    "regexparam": "^1.3.0",
    "svelte": "^3.23.2"
  }
}
`),
	"/plenti.json": []byte(`{
	"types": {
		"pages": "/:filename"
	},
	"build": "public",
	"local": {
		"port": 3000
	}
}`),
}
