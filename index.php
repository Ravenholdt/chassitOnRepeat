<html>

<head>
    <?php
        $name = $_GET["v"];
        $start = $_GET["s"];
        $end = $_GET["e"];

        $file = glob('files/*-'. $name .'.mp4');

        include "dbconnect.php";
    ?>
    <link rel="stylesheet" href="index.css">
</head>

<body>
    <div id="video">
        <video id="myvideo" controls autoplay>
            <source src="<?php echo $file[0]; ?>" type="video/mp4">
            Your browser does not support HTML5 video.
        </video>

        <?php
        if (!isset($_GET["s"])){ $start = "0"; }
        if (!isset($_GET["e"])){ $end = "1000000"; }
        ?>
    </div>

    <div id="history">
        <?php include "history.php"; ?>
    </div>

    <script>
        myVideo = document.getElementById("myvideo");

        let start = <?php echo $start ?>;
        let end = <?php echo $end ?>;

        console.log(<?php echo "\"Total playtime: " . $totalTime . "\""; ?>);
        console.log(<?php echo "\"Total playtime: " . $totalTime/(3600*24) . " Days\""; ?>);
        console.log(start);
        console.log(end);

        /*myVideo.onloadedmetadata = function() {
            console.log('metadata loaded!');
            console.log(this.duration);//this refers to myVideo
            if (end > 999999){
                end = this.duration;
            }
        };*/

        //var deltaTime = 0;
        //var lastTime = myVideo.currentTime;

        /*
        document.getElementById('myvideo').addEventListener('loadedmetadata', function() {
            this.currentTime = start; 
        }, false);
        */

        //echo myVideo.duration;

        

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

        //setTimeout(update(deltaTime), 600000);

        myVideo.addEventListener('timeupdate', function() {
            console.log(this.currentTime);
            //console.log(deltaTime);

            if (this.currentTime < start) {
                this.currentTime = start;
                //lastTime = this.currentTime;
                console.log("Jump Forward");
            }

            //deltaTime += this.currentTime - lastTime;
            //lastTime = this.currentTime;

            if (this.currentTime > end) {
                this.currentTime = start;
                //lastTime = this.currentTime;
                console.log("Restart");
                myVideo.play();
                update(end-start);
            }
        }, false);

        myVideo.addEventListener('ended', function() {
            this.currentTime = start;
            //lastTime = this.currentTime;
            myVideo.play();
            console.log("Ended");
            update(this.duration-start);
        },false);

        myVideo.play();
    </script>

    <br>
</body>

</html>