<?php

use Chassit\Repeat\History;
use Chassit\Repeat\Video;

require_once "init.php";



/*
if (!isset($_GET["s"])){ $start = "0"; }
if (!isset($_GET["e"])){ $end = "1000000"; }
*/

$video = History::getRandomVideo();
$file = "files/$video->name-$video->id.mp4";

$start = $video->start ?? 0;
$end = $video->end ?? 100000000;

$totalTime = History::getTotalTime();
?>
<!DOCTYPE html>
<html lang="sv">
<head>
    <link rel="stylesheet" href="index.css">
    <title>(0)Chassit radio - <?= $video->name ?></title>
</head>

<body>
    <video id="myvideo" controls autoplay style="width:100%;max-height: 800px;">
        <source src="<?= $file ?>" type="video/mp4">
        Your browser does not support HTML5 video.
    </video>
    <h4 id="videoTitle"> <?= $video->name ?> </h4>
    <p id="videoLoops"> <?= History::toDisplayTime($video->playtime, true) ?> </p>
    <p id="switches"> 0 videos resisted </p>

<script>
    const myVideo = document.getElementById("myvideo");
    myVideo.volume = 0.5;
    const videoTitle = document.getElementById("videoTitle");
    const videoLoops = document.getElementById("videoLoops");
    const videoSwitches = document.getElementById("switches");

    let start = <?= $start ?>;
    let end = <?= $end ?>;
    let id = '<?= $video->id ?>';
    let switches = 0;
    console.log(<?php echo "\"Total playtime: " . $totalTime . "s\""; ?>);
    console.log(<?php echo "\"Total playtime: " . History::toDisplayTime($totalTime, true) . "\""; ?>);
    console.log(start);
    console.log(end);

    function update(t) {
        fetch("../update.php", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                v: id,
                t: t,
                s: start === 0 ? null: start,
                e: end === 100000000 ? null : end,
            })
        }).then(async value => {
            if (!value.ok)
                console.log(value.status, value.statusText, await value.text());
        });
    }

    function endHandler(){
        update(this.duration - start);
        fetch("/switchrandom.php").then(e => e.json()).then(e => {
            console.log(e);
            start = e.start ?? 0;
            end = e.end ?? 100000000;
            id = e.id;
            myVideo.src = `files/${e.name}-${e.id}.mp4`;
            myVideo.currentTime = start;
            console.log(`${start}, ${end}, ${id}, ${myVideo.currentTime} & ${myVideo.src}`);
            videoTitle.innerText = e.name;
            videoLoops.innerText = `Loops: ${getTimeStr(e.playtime)}`;
            switches++;
            videoSwitches.innerText = `${switches} videos resisted`;
            document.title = `(${switches})Chassit radio - ${e.name}`;
            myVideo.play();
            console.log("Switched")
        }).catch(err => console.error(err))
        console.log("Ended");
    }

    function getTimeStr(time){
        let days = Math.floor(time / 86400);
        time -= days * 86400;
        let hours = Math.floor(time / 3600);
        time -= hours * 3600;
        let minutes = Math.floor(time / 60);
        time -= minutes * 60;
        let seconds = Math.floor(time % 60);
        return `${days}d ${hours}h ${minutes}m ${seconds}s`;
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

</script>

<br>
</body>

</html>
