<script>
    export let mediaList, buttonText, action, encoding;
    import { publish } from './publish.js';

    let status;
    const onSubmit = async () => {
        status = "sending";
        try {
            await publish(mediaList, action, encoding);
            status = "sent";
            resetStatus();
        } catch (error) {
            status = "failed";
            resetStatus();
            throw error;
        }
    }
    const resetStatus = () => {
        setTimeout(() => {
            status = "";
            mediaList = [];
        }, 700);
    }
</script>

<button 
    on:click|preventDefault={onSubmit}
    type="submit"
    disabled={status}
>
    {#if status == "sending"}
        Sending...
    {:else if status == "failed"}
        Could not commit the changes.
    {:else if status == "sent"}
        Changes committed.
    {:else}
        {buttonText}
    {/if}
</button>