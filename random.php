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

$start = isset($video->start) ? $video->start : 0;
$end = isset($video->end) ? $video->end : 100000000;

$totalTime = History::getTotalTime();
?>
<!DOCTYPE html>
<html lang="sv">
<head>
    <link rel="stylesheet" href="index.css">
    <title><?= $video->name ?> - Chassit on Repeat</title>
</head>

<body>
    <video id="myvideo" controls autoplay style="width:100%;">
        <source src="<?= $file ?>" type="video/mp4">
        Your browser does not support HTML5 video.
    </video>

<script>
    myVideo = document.getElementById("myvideo");
    myVideo.volume = 0.5;

    let start = <?= $start ?>;
    let end = <?= $end ?>;
    let id = '<?= $video->id ?>';
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
                s: start,
                e: end,
            })
        }).then(async value => {
            if (!value.ok)
                console.log(value.status, value.statusText, await value.text());
        });
    }

    function endHandler(){
        update(this.duration - start);
        fetch("/switchrandom.php").then(e => e.json()).then(e => {
            start = e.start ?? 0;
            end = e.end ?? 100000000;
            id = e.id;
            myVideo.src = `files/${e.name}-${e.id}.mp4`;
            this.currentTime = start;
            myVideo.play();
            console.log("Switched")
        }).catch(err => console.error(err))
        console.log("Ended");
        myVideo.play();
    }

    myVideo.addEventListener('timeupdate', function () {

        if (this.currentTime < start) {
            this.currentTime = start;
            console.log("Jump Forward");
        }

        if (this.currentTime > end) {
            endHandler();
            console.log("Restart");
            
        }
    }, false);

    myVideo.addEventListener('ended', endHandler, false);
        
</script>

<br>
</body>

</html>