<script lang="ts">
    import Button from "./lib/Button.svelte";
    import QuizCard from "./lib/QuizCard.svelte";

    let quizzes: { id: string; name: string }[] = [];

    let code = "";
    let msg = ""

    async function getQuizzes() {
        let response = await fetch("http://localhost:3000/api/quizzes");
        if (!response.ok) {
            console.error("Error:", response.statusText); // Muestra el mensaje de error
            alert("Falló");
            return;
        }
        let json = await response.json();
        quizzes = json;
        console.log(json);
    }

    function connect() {
        let websocket = new WebSocket("ws://localhost:3000/ws");
        websocket.onopen = () => {
            console.log("opened connection");
            websocket.send(`join: ${code}`);
        };

        websocket.onmessage = (event) => {
            console.log(event.data);
        };
    }

    function hostQuiz(quiz: { id: string; name: string }) {
        let websocket = new WebSocket("ws://localhost:3000/ws");
        websocket.onopen = () => {
            console.log("opened connection");
            websocket.send(`host: ${quiz.id}`);
        };

        websocket.onmessage = (event) => {
            msg = event.data;
            console.log(event.data);
        };
    }
</script>

<main class="p-2 grid place-content-center h-screen gap-4">
    <!-- <Button>Cool</Button> -->

    <button
        class="px-3 py-1 shadow-lg shadow-gray-500/50 bg-black text-white rounded-lg text-[15px] cursor-pointer active:scale-[.97]"
        on:click={getQuizzes}>Get Quizzes</button
    >
    Message: {msg}

    {#each quizzes as quiz}
        <QuizCard {quiz} host={hostQuiz} />
    {/each}
    <input bind:value={code} type="text" class="border" placeholder="Game code" />
    <button
        class="px-3 py-1 shadow-lg shadow-gray-500/50 bg-black text-white rounded-lg text-[15px] cursor-pointer active:scale-[.97]"
        on:click={connect}>Join game</button
    >
</main>
