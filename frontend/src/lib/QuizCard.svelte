<!-- QuizCard.svelte -->
<script lang="ts">
    import Button from "./Button.svelte";

    export let quiz: { id: string; name: string };
    export let host: (quiz: { id: string; name: string }) => void;
</script>

<div class="bg-white border-2 border-green-500 p-4 rounded-xl">
    <p>{quiz.name}</p>
    <Button on:click={() => host(quiz)}>Host</Button>
</div>
