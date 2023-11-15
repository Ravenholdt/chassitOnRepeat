(function (){
    const data = elementById("video-data")
    let id = data.dataset.id;
    let start = parseFloat(data.dataset.start);
    let end = parseFloat(data.dataset.end);
    let safe = data.dataset.safe === 'true';
    data.remove();

    let switches = 0;
    console.log(`Start: ${start}`);
    console.log(`End: ${end}`);

    /**
     * @type {HTMLVideoElement}
     */
    const videoElement = elementById("my-video");
    videoElement.volume = 0.5;

    let updateInterval = setInterval(updateCheckFunc(videoElement, start, end, endHandler), 40);

    const videoTitle = elementById("video-title");
    const videoLoops = elementById("video-loops");
    const totalLoops = elementById("total-loops");
    const videoSwitches = elementById("switches");

    /**
     * Updates the played time on the playlist
     * @param {number} t The played time
     * @returns {Promise<void>}
     */
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
            totalLoops.innerText = json.time_formatted;
        } else
            console.log(value.status, value.statusText, await value.text());
    }

    /**
     * Updates the video with a new one from the api
     * @returns {Promise<void>}
     */
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
        videoLoops.innerText = json.time_formatted;

        switches++;
        videoSwitches.innerText = `${switches}`;
        document.title = `(${switches}) Chassit radio - ${json.title}`;

        await videoElement.play();
        updateInterval = setInterval(updateCheckFunc(videoElement, start, end, endHandler), 40);
        console.log("Switched")
    }

    /**
     * Handler for when the video ends
     * @param {number} t The played time
     * @returns {Promise<void>}
     */
    async function endHandler(t){
        clearInterval(updateInterval);

        await Promise.all([update(t), updateVideo()]);
        console.log("Ended");
    }
})();