(function (){
    const data = document.getElementById("video-data")
    let id = data.dataset.id;
    let start = parseFloat(data.dataset.start);
    let end = parseFloat(data.dataset.end);
    let safe = data.dataset.safe === 'safe';
    let totalPlaytime = parseInt(data.dataset.total_play_time, 10);
    let formattedTotalPlaytime =  data.dataset.formatted_total_play_time;
    data.remove();

    if (typeof id === 'undefined')
    {
        console.log("No id set, ignoring")
        return;
    }

    console.log("Total playtime: " + totalPlaytime);
    console.log("Total playtime: " + formattedTotalPlaytime);
    console.log("Start: " + start);
    console.log("End: " + end);

    const myVideo = document.getElementById("my-video");
    myVideo.volume = 0.5;

    const videoLoops = document.getElementById("video-loops");

    function update(t) {
        fetch(`/api/v1/video/${id}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                time: t,
            })
        }).then(async value => {
            if (!value.ok)
                console.log(value.status, value.statusText, await value.text());
            else {
                let json = await value.json();
                videoLoops.innerText = "Playtime: " + json.time_formatted;
            }
        });
    }

    function sendInterval() {
        fetch(`/api/v1/video/${id}/settings`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                start: start > 0 ? start : null,
                end: end > myVideo.duration ? null : end,
                safe: safe
            })
        }).then(async value => {
            if (value.ok) {
                const success = document.getElementById("update-success");

                success.classList.remove("o-0")
                setTimeout(() => {
                    success.classList.add("o-0");
                }, 2000);
            } else {
                console.log(value.status, value.statusText, await value.text());
            }
        });
    }

    document.getElementById("update-loop").addEventListener("click", (event) => {
        sendInterval();
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


    myVideo.addEventListener('timeupdate', function () {
        if (this.currentTime < start) {
            this.currentTime = start;
            console.log("Jump Forward");
        }

        if (this.currentTime > end) {
            this.currentTime = start;
            console.log("Restart");
            myVideo.play();
            update(end - start);
        }
    }, false);


    myVideo.addEventListener('ended', function () {
        this.currentTime = start;
        myVideo.play();
        console.log("Ended");
        update(this.duration - start);
    }, false);


    myVideo.play();
})();