<script lang="ts">
    import Button from "./lib/Button.svelte";
    import QuizCard from "./lib/QuizCard.svelte";
    import type { Player, QuizQuestion } from "./model/quiz";
    import { NetService, PacketTypes, type ChangeGameStatePacket, type PlayerJoinPacket, type TickPacket } from "./service/net";

    let quizzes: { id: string; name: string }[] = [];
    let currentQuestion: QuizQuestion | null = null;

    let tick = 0;

    let players: Player[] = [];

    let netService = new NetService();

    setTimeout(() => {
        netService.connect();
    },500)

    netService.onPacket((packet: any) => {
        console.log(packet);
        switch (packet.id) {
            case 2: {
                currentQuestion = packet.question;
                break;
            }

            case PacketTypes.ChangeGameState: {
                let data = packet as ChangeGameStatePacket;
                console.log(data.state);
                break;
            }
            case PacketTypes.PlayerJoin:{
                let data = packet as PlayerJoinPacket;
                players = [...players,data.player];
                break;
            }
            case PacketTypes.Tick: {
                let data = packet as TickPacket;
                tick = data.tick;
                break
            }
        }
    });

    let code = "";
    let name = "";
    let msg = "";

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
        netService.sendPacket({
            id: 0,
            code: code,
            name: name,
        });

        // let websocket = new WebSocket("ws://localhost:3000/ws");
        // websocket.onopen = () => {
        //     console.log("opened connection");
        //     websocket.send(`join: ${code}`);
        // };

        // websocket.onmessage = (event) => {
        //     console.log(event.data);
        // };
    }

    function startGame(){
        netService.sendPacket({
            id: PacketTypes.StartGame,
        });
    }

    function hostQuiz(quiz: { id: string; name: string }) {
        netService.sendPacket({
            id: 1,
            quizId: quiz.id,
        });
    }
</script>

<main
    class="p-2 grid grid-cols-3 grid-rows-4 h-screen gap-4"
>
    <!-- <Button>Cool</Button> -->

    <div class="col-start-2 flex flex-col gap-4">
        <button
            class="px-3 py-3 shadow-lg shadow-gray-500/50 bg-black text-white rounded-lg text-[15px] cursor-pointer active:scale-[.97]"
            on:click={getQuizzes}>Get Quizzes</button
        >

        Message: {msg}

        {#each quizzes as quiz}
            <QuizCard {quiz} host={hostQuiz} />
        {/each}
        <input
            bind:value={code}
            type="text"
            class="border"
            placeholder="Game code"
        />
        <input
            bind:value={name}
            type="text"
            class="border"
            placeholder="Name"
        />
        <button
            class="px-3 py-3 shadow-lg shadow-gray-500/50 bg-black text-white rounded-lg text-[15px] cursor-pointer active:scale-[.97]"
            on:click={connect}>Join game</button
        >
    </div>
    <div class="row-start-3 col-span-3">
        {#if currentQuestion != null}
            <h2 class="text-4xl font-bold mt-8">{currentQuestion.name}</h2>
            <div class="flex">
                {#each currentQuestion.choices as choice}
                    <div
                        class="flex-1 bg-blue-400 text-center font-bold text-2"
                    >
                        {choice.name}
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</main>
