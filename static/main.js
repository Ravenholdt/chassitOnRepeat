(function (){
    const data = elementById("video-data")
    let id = data.dataset.id;
    let start = parseFloat(data.dataset.start);
    let end = parseFloat(data.dataset.end);
    let safe = data.dataset.safe === 'true';
    let totalPlaytime = parseInt(data.dataset.total_play_time, 10);
    let formattedTotalPlaytime =  data.dataset.formatted_total_play_time;
    data.remove();

    if (typeof id === 'undefined') {
        console.log("No id set, ignoring")
        return;
    }

    console.log(`Total playtime: ${totalPlaytime}`);
    console.log(`Total playtime: ${formattedTotalPlaytime}`);
    console.log(`Start: ${start}`);
    console.log(`End: ${end}`);
    console.log(`Safe: ${safe}`)

    /**
     * @type {HTMLVideoElement}
     */
    const videoElement = elementById("my-video");
    videoElement.volume = 0.5;

    const videoLoops = elementById("video-loops");
    const currentLoops = elementById("current-loops");
    let currentTime = 0;

    let updateTimeInterval;
    function stopTimer() {
        clearInterval(updateTimeInterval);
    }

    function startTimer() {
        updateTimeInterval = setInterval(updateTime, 40);
    }

    /**
     * Update the current playtime
     * @param {number} time
     */
    function updateCurrentPlaytime(time) {
        currentTime += time;
        currentLoops.innerText = formatTime(currentTime);
    }

    /**
     * Updates the current time played for the video
     * @param {number} t The played time to add to the video
     */
    async function update(t) {
        stopTimer();

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
            videoLoops.innerText = json.time_formatted;
            updateCurrentPlaytime(Math.floor(t));
        } else
            console.log(value.status, value.statusText, await value.text());

        startTimer();
    }

    async function sendInterval() {
        const value = await fetch(`/api/v1/video/${id}/settings`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                start: start > 0 ? start : null,
                end: end > videoElement.duration ? null : end,
                safe: safe
            })
        });
        if (value.ok) {
            const success = elementById("update-success");

            success.classList.remove("o-0")
            setTimeout(() => {
                success.classList.add("o-0");
            }, 2000);
        } else
            console.log(value.status, value.statusText, await value.text());
    }

    elementById("update-loop").addEventListener("click", async (event) => {
        await sendInterval();
    });


    elementById("safe").addEventListener("change",(event) =>{
        safe = event.target.checked;
    });


    elementById("start").addEventListener("change",(event) =>{
        let s = event.target.value;
        if (s === "")
            start = 0;
        else {
            s = parseFloat(s);
            if(isNaN(s)) {
                return;
            }
            start = s;
        }
    });

    elementById("end").addEventListener("change", (event) => {
        let e = event.target.value;
        if (e === "")
            end = 90000;
        else {
            e = parseFloat(e);
            if(isNaN(e)) {
                return;
            }
            end = e;
        }
    });

    videoElement.addEventListener("play", function () {
        // Use an internal timer instead of video as this is more accurate and fires more often
        startTimer();
    })

    videoElement.addEventListener("pause", async function () {
        stopTimer();
    })

    async function updateTime() {
        if (videoElement.currentTime < start) {
            videoElement.currentTime = start;
            console.log("Jump Forward");
            return;
        }

        if (videoElement.currentTime >= (videoElement.duration - 0.05)) {
            // User browser loop when the start isn't set
            if (start > 0)
                videoElement.currentTime = start;

            console.log("Ended")
            await update(videoElement.duration - start);
            return;
        }

        if (videoElement.currentTime >= end) {
            videoElement.currentTime = start;
            await videoElement.play()

            console.log("Restart");
            await update(end - start);
        }
    }

    // Try to autostart the video
    videoElement.play().catch(reason => {
        console.log("Error starting video automatically", reason);
    });
})();