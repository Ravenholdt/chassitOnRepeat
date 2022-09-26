<?php
    $name = "";
    $start = "0";
    $end = "1000000";

    if (isset($_GET["v"])){ $name = $_GET["v"]; }
    if (isset($_GET["s"])){ $start = $_GET["s"]; }
    if (isset($_GET["e"])){ $end = $_GET["e"]; }

    $file = glob('files/*-'. $name .'.mp4');

    include "dbconnect.php";
?>

<html>

<head>
    <link rel="stylesheet" href="index.css">
    <title><?php echo explode(".mp4", explode("/", $file[0])[1])[0] ?> - Chassit on Repeat</title>
</head>

<body>
    <div id="video">
        <video id="myvideo" controls autoplay>
            <source src="<?php echo $file[0]; ?>" type="video/mp4">
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
        <?php include "history.php"; ?>
    </div>

    <script>
        myVideo = document.getElementById("myvideo");
        myVideo.volume = 0.5;

        let start = <?php echo $start ?>;
        let end = <?php echo $end ?>;

        console.log(<?php echo "\"Total playtime: " . $totalTime . "\""; ?>);
        console.log(<?php echo "\"Total playtime: " . $totalTime/(3600*24) . " Days\""; ?>);
        console.log(start);
        console.log(end);        

        function update(t) {
            deltaTime = 0;
            let xmlHttp = new XMLHttpRequest();

            <?php
                $updateStart = "";
                $updateEnd = "";

                if (isset($_GET["s"])){ $updateStart = "+\"&s=" . $_GET["s"] . "\""; }
                if (isset($_GET["e"])){ $updateEnd = "+\"&e=" . $_GET["e"] . "\""; }
            ?>

            let url="update.php?v=<?php echo $name; ?>&t="+t <?php echo $updateStart . $updateEnd; ?> ;
            xmlHttp.open( "GET", url, false ); // false for synchronous request
            xmlHttp.send();
            console.log(xmlHttp.responseText);
        }

        myVideo.addEventListener('timeupdate', function() {

            if (this.currentTime < start) {
                this.currentTime = start;
                console.log("Jump Forward");
            }

            if (this.currentTime > end) {
                this.currentTime = start;
                console.log("Restart");
                myVideo.play();
                update(end-start);
            }
        }, false);

        myVideo.addEventListener('ended', function() {
            this.currentTime = start;
            myVideo.play();
            console.log("Ended");
            update(this.duration-start);
        },false);

        myVideo.play();
    </script>

    <br>
</body>

</html>