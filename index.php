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

$id = "";
$start = "0";
$end = "1000000";

if (isset($_GET["v"])) {
    $id = $_GET["v"];
}
if (isset($_GET["s"])) {
    $start = $_GET["s"];
}
if (isset($_GET["e"])) {
    $end = $_GET["e"];
}

$video = getVideo($id);
$file = "files/$video->name-$video->id.mp4"
?>

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

    <?php
    /*
    if (!isset($_GET["s"])){ $start = "0"; }
    if (!isset($_GET["e"])){ $end = "1000000"; }
    */
    ?>
</div>

<div id="history">
    <?php
    History::loadVideos();
    History::render();
    $totalTime = History::getTotalTime();
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
        <?php
        $args = "";

        if (isset($_GET["s"])) {
            $args .= "&s={$_GET["s"]}";
        }
        if (isset($_GET["e"])) {
            $args .= "&e={$_GET["e"]}";
        }
        ?>
        fetch("update.php", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                v: '<?= $id ?>',
                t: t,
                s: '<?= $_GET["s"] ?? null ?>',
                e: '<?= $_GET["e"] ?? null ?>',
            })
        }).then(value => console.log(value));
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