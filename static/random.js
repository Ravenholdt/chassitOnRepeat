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

    const myVideo = document.getElementById("my-video");
    myVideo.volume = 0.5;

    const videoTitle = document.getElementById("video-title");
    const videoLoops = document.getElementById("video-loops");
    const videoSwitches = document.getElementById("switches");

    function endHandler(){
        fetch("/api/v1/video/random" + (safe ? '?safe': ''))
            .then(async resp => {
                const json = await resp.json();
                console.log(json);

                start = json.start;
                end = json.end;
                id = json.id;

                myVideo.src = json.url;
                myVideo.currentTime = start;

                console.log(`${start}, ${end}, ${id}, ${myVideo.currentTime} & ${myVideo.src}`);

                videoTitle.innerText = json.title;
                videoLoops.innerText = `Playtime: ${json.time_formatted}`;

                switches++;
                videoSwitches.innerText = `${switches} videos resisted`;
                document.title = `(${switches})Chassit radio - ${json.title}`;

                myVideo.play();
                console.log("Switched")
            })
            .catch(err => console.error(err))

        console.log("Ended");
    }

    myVideo.addEventListener('timeupdate', function () {

        if (this.currentTime < start) {
            this.currentTime = start;
            console.log("Jump Forward");
        }

        if (this.currentTime > end) {
            console.log("Restart");
            endHandler();
        }
    }, false);

    myVideo.addEventListener('ended', endHandler, false);
})();