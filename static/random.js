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
    const videoSwitches = document.getElementById("switches");

    function endHandler(){
        clearInterval(updateInterval)
        fetch("/api/v1/video/random" + (safe ? '?safe': ''))
            .then(async resp => {
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
            })
            .catch(err => console.error(err))

        console.log("Ended");
    }

    function updateTimer() {
        if (videoElement.currentTime < start) {
            videoElement.currentTime = start;
            console.log("Jump Forward");
        }

        if (videoElement.currentTime >= end) {
            console.log("Restart");
            endHandler();
        }

        if (videoElement.currentTime >= videoElement.duration) {
            endHandler();
        }
    }
})();