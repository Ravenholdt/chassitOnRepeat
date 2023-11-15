(function (){
    const data = document.getElementById("video-data")
    let start = parseFloat(data.dataset.start);
    let end = parseFloat(data.dataset.end);
    let playlist_id = data.dataset.playlist_id;
    let playlist_name = data.dataset.playlist_name;
    data.remove();

    console.log(start);
    console.log(end);

    const videoElement = document.getElementById("my-video");
    videoElement.volume = 0.5;

    let updateInterval = setInterval(updateTimer, 40);

    const videoTitle = document.getElementById("video-title");
    const videoLoops = document.getElementById("video-loops");
    const totalLoops = document.getElementById("total-loops");
    const currentLoops = document.getElementById("current-loops");

    let currentTime = 0;

    async function update(t) {
        const value = await fetch(`/api/v1/playlist/${playlist_id}`, {
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
            totalLoops.innerText = json.time_formatted;
        } else
            console.log(value.status, value.statusText, await value.text());
    }

    async function updateVideo() {
        const resp = await fetch(`/api/v1/playlist/${playlist_id}/random`)
        const json = await resp.json();

        start = json.start;
        end = json.end;

        videoElement.src = json.url;
        videoElement.currentTime = start;

        console.log(start, end, videoElement.currentTime, videoElement.src);

        videoTitle.innerText = json.title;
        videoLoops.innerText = json.time_formatted;

        document.title = `${json.title} - ${playlist_name} - Chassit on Repeat`;

        videoElement.play();
        updateInterval = setInterval(updateTimer, 40);
        console.log("Switched")
    }

    async function endHandler(t){
        clearInterval(updateInterval);

        currentTime += t;
        currentLoops.innerText = formatTime(currentTime)

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