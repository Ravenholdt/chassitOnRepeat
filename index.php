<?php

use Chassit\Repeat\History;
use Chassit\Repeat\Video;

require_once "init.php";

function getVideo($id): Video
{
    foreach (glob("files/*-$id.mp4") as $filename) {
        preg_match("/^files\/(.*)-([A-Za-z0-9_-]{11}).mp4$/", $filename, $matches);
        return new Video($matches[2], $matches[1]);
    }
    return new Video("", "");
}

function getFloatParam($key, ?float $defaultVal = null): ?float
{
    if (isset($_GET[$key]) && is_numeric($_GET[$key]))
        return floatval($_GET[$key]);
    return $defaultVal;
}

$id = "";
$start = getFloatParam('s', 0);
$end = getFloatParam('e', 1000000);

if (isset($_GET["v"])) {
    $id = $_GET["v"];
}

/*
if (!isset($_GET["s"])){ $start = "0"; }
if (!isset($_GET["e"])){ $end = "1000000"; }
*/

$video = getVideo($id);
$file = "files/$video->name-$video->id.mp4";

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

<div id="history">
    <?php
    History::render();
    ?>
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
        fetch("update.php", {
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
        this.currentTime = start;
        myVideo.play();
        console.log("Ended");
        update(this.duration - start);
    }, false);

    myVideo.play();
</script>

<br>
</body>

</html>