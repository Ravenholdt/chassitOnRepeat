(function (){
    const data = document.getElementById("video-data")
    let id = data.dataset.id;
    let start = parseFloat(data.dataset.start);
    let end = parseFloat(data.dataset.end);
    let safe = data.dataset.safe === 'true';
    data.remove();

    let switches = 0;
    console.log(start);
    console.log(end);

    const videoElement = document.getElementById("my-video");
    videoElement.volume = 0.5;

    let updateInterval = setInterval(updateTimer, 40);

    const videoTitle = document.getElementById("video-title");
    const videoLoops = document.getElementById("video-loops");
    const totalLoops = document.getElementById("total-loops");
    const videoSwitches = document.getElementById("switches");

    async function update(t) {

        const id = safe ? 'RANDOM-SAFE': 'RANDOM';
        const value = await fetch(`/api/v1/video/${id}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                time: t,
            })
        });

        if (value.ok) {
            let json = await value.json();
            totalLoops.innerText = "Total Random Playtime: " + json.time_formatted;
        } else
            console.log(value.status, value.statusText, await value.text());
    }

    async function updateVideo() {
        const resp = await fetch("/api/v1/video/random" + (safe ? '?safe': ''))
        const json = await resp.json();
        console.log(json);

        start = json.start;
        end = json.end;
        id = json.id;

        videoElement.src = json.url;
        videoElement.currentTime = start;

        console.log(start, end, id, videoElement.currentTime, videoElement.src);

        videoTitle.innerText = json.title;
        videoLoops.innerText = `Playtime: ${json.time_formatted}`;

        switches++;
        videoSwitches.innerText = `${switches} videos resisted`;
        document.title = `(${switches}) Chassit radio - ${json.title}`;

        videoElement.play();
        updateInterval = setInterval(updateTimer, 40);
        console.log("Switched")
    }

    async function endHandler(t){
        clearInterval(updateInterval);

        await Promise.all([update(t), updateVideo()]);
        console.log("Ended");
    }

    async function updateTimer() {
        if (videoElement.currentTime < start) {
            videoElement.currentTime = start;
            console.log("Jump Forward");
        }

        if (videoElement.currentTime >= end) {
            console.log("Restart");
            await endHandler(end - start);
        }

        if (videoElement.currentTime >= videoElement.duration) {
            await endHandler(videoElement.duration - start);
        }
    }
})();