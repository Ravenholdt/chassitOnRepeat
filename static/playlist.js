(function (){
    const data = elementById("video-data")
    let start = parseFloat(data.dataset.start);
    let end = parseFloat(data.dataset.end);
    let playlist_id = data.dataset.playlist_id;
    let playlist_name = data.dataset.playlist_name;
    data.remove();

    /**
     * @type {HTMLVideoElement}
     */
    const videoElement = elementById("my-video");
    videoElement.volume = 0.5;

    console.log(start, end, videoElement.currentTime, videoElement.src);

    let updateInterval = setInterval(updateCheckFunc(videoElement, start, end, endHandler), 40);

    const videoTitle = elementById("video-title");
    const videoLoops = elementById("video-loops");
    const totalLoops = elementById("total-loops");
    const currentLoops = elementById("current-loops");

    let currentTime = 0;

    /**
     * Updates the played time on the playlist
     * @param {number} t The played time
     * @returns {Promise<void>}
     */
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

    /**
     * Updates the video with a new one from the api
     * @returns {Promise<void>}
     */
    async function updateVideo() {
        await updateVideoInfo(`/api/v1/playlist/${playlist_id}/random`)
    }

    /**
     * Updates the page with video data from a json response
     * @param {string} url The url that will fetch a new video
     * @returns {Promise<void>}
     */
    async function updateVideoInfo(url) {
        const resp = await fetch(url)
        const json = await resp.json();

        start = json.start;
        end = json.end;
        videoElement.src = json.url;
        videoElement.currentTime = start;

        console.log(start, end, videoElement.currentTime, videoElement.src);

        videoTitle.innerText = json.title;
        videoLoops.innerText = json.time_formatted;

        document.title = `${json.title} - ${playlist_name} - Chassit on Repeat`;
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

        currentTime += t;
        currentLoops.innerText = formatTime(currentTime)

        await Promise.all([update(t), updateVideo()]);
        console.log("Ended");
    }

    /**
     * Changes the video to the clicked video
     * @param {MouseEvent} e The click event
     */
    async function videoClicked(e) {
        await updateVideoInfo(`/api/v1/video/${this.dataset.id}`)
    }

    // Adds click event listeners to all playlist videos
    Array.from(elementsByClass("playlist-video")).forEach(function(e) {
        e.addEventListener("click", videoClicked)
    })
})();