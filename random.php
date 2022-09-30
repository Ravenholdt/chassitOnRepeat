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

$start = $video->$start ?? 0;
$end = $video->$end ?? 100000000;

$totalTime = History::getTotalTime();
?>
<!DOCTYPE html>
<html lang="sv">
<head>
    <link rel="stylesheet" href="index.css">
    <title><?= $video->name ?> - Chassit on Repeat</title>
</head>

<body>
<div id="video">
    <video id="myvideo" controls autoplay>
        <source src="<?= $file ?>" type="video/mp4">
        Your browser does not support HTML5 video.
    </video>
</div>

<script>
    myVideo = document.getElementById("myvideo");
    myVideo.volume = 0.5;

    let start = <?= $start ?>;
    let end = <?= $end ?>;
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
                v: '<?= $id ?>',
                t: t,
                s: <?= getFloatParam('s') ?? 'null' ?>,
                e: <?= getFloatParam('e') ?? 'null' ?>,
            })
        }).then(async value => {
            if (!value.ok)
                console.log(value.status, value.statusText, await value.text());
        });
    }

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
        update(this.duration - start);
        fetch("/switch.php").then(e => e.json()).then(e => {
            start = e.start ?? 0;
            end = e.end ?? 100000000;
            myVideo.src = `files/${e.name}-${e.id}.mp4`;
            this.currentTime = start;
            myVideo.play();
            console.log("Switched")
        }).catch(err => console.error(err))
        console.log("Ended");
    }, false);

    myVideo.play();
</script>

<br>
</body>

</html>