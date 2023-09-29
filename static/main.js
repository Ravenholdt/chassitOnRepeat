(function (){
    const data = document.getElementById("video-data")
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

    console.log("Total playtime: " + totalPlaytime);
    console.log("Total playtime: " + formattedTotalPlaytime);
    console.log("Start: " + start);
    console.log("End: " + end);
    console.log("Safe: " + safe)

    const videoElement = document.getElementById("my-video");
    videoElement.volume = 0.5;

    const videoLoops = document.getElementById("video-loops");
    const currentLoops = document.getElementById("current-loops");
    let currentTime = 0;

    let updateTimeInterval;
    function stopTimer() {
        clearInterval(updateTimeInterval);
    }

    function startTimer() {
        updateTimeInterval = setInterval(updateTime, 40);
    }

    /**
     * Format a number to the given amount of digits
     * @param {number} num
     * @param {number} digits
     */
    function formatNumber(num, digits = 2) {
        return num.toLocaleString('en-US', {minimumIntegerDigits: digits, useGrouping: false})
    }

    /**
     * Formats the time to show days, hours, minutes and seconds
     * @param {number} time
     */
    function formatTime(time) {
        const hourSeconds = time % (60 * 60 * 24);
        const minuteSeconds = hourSeconds % (60 * 60)
        const remainingSeconds = minuteSeconds % 60

        const days = Math.floor(time / (60 * 60 * 24))
        const hours = Math.floor(hourSeconds / (60 * 60))
        const minutes = Math.floor(minuteSeconds / 60)
        const seconds = Math.floor(remainingSeconds)

        if (days > 0)
            return `${days}d ${formatNumber(hours)}h ${formatNumber(minutes)}m ${formatNumber(seconds)}s`;
        if (hours > 0)
            return `${formatNumber(hours)}h ${formatNumber(minutes)}m ${formatNumber(seconds)}s`
        if (minutes > 0)
            return `${formatNumber(minutes)}m ${formatNumber(seconds)}s`
        return `${seconds}s`
    }

    /**
     * Update the current playtime
     * @param {number} time
     */
    function updateCurrentPlaytime(time) {
        currentTime += time;
        currentLoops.innerText = "Current playtime: " + formatTime(currentTime);
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
            videoLoops.innerText = "Total playtime: " + json.time_formatted;
            updateCurrentPlaytime(t);
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
            const success = document.getElementById("update-success");

            success.classList.remove("o-0")
            setTimeout(() => {
                success.classList.add("o-0");
            }, 2000);
        } else
            console.log(value.status, value.statusText, await value.text());
    }

    document.getElementById("update-loop").addEventListener("click", async (event) => {
        await sendInterval();
    });


    document.getElementById("safe").addEventListener("change",(event) =>{
        safe = event.target.checked;
    });


    document.getElementById("start").addEventListener("change",(event) =>{
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

    document.getElementById("end").addEventListener("change", (event) => {
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
            videoElement.play()

            console.log("Restart");
            await update(end - start);
        }
    }

    // Try to autostart the video
    videoElement.play().catch(reason => {
        console.log("Error starting video automatically", reason);
    });
})();