<html>

<?php
    $name = $_GET["v"];
    $start = $_GET["s"];
    $end = $_GET["e"];

    $file = glob('files/*-'. $name .'.mp4');
?>

<header>
</header>

<body>
    <video id="myvideo" controls autoplay>
        <source src="<?php echo $file[0]; ?>" type="video/mp4">
        Your browser does not support HTML5 video.
    </video>

    <script>
        myVideo = document.getElementById("myvideo");

        <?php
        if (!isset($_GET["s"])){ $start = 0; }
        if (!isset($_GET["e"])){ $end = -1; }
        ?>

        var start = <?php echo (int)$start ?>;
        var end = <?php echo (int)$end ?>;

        myVideo.onloadedmetadata = function() {
            console.log('metadata loaded!');
            console.log(this.duration);//this refers to myVideo
            if (end < 0){
                end = this.duration;
            }
        };

        var deltaTime = 0;
        var lastTime = myVideo.currentTime;

        /*
        document.getElementById('myvideo').addEventListener('loadedmetadata', function() {
            this.currentTime = start; 
        }, false);
        */

        //echo myVideo.duration;

        

        function update(t) {
            deltaTime = 0;
            var xmlHttp = new XMLHttpRequest();

            var url="update.php?v=<?php echo $name; ?>&t="+t;
            xmlHttp.open( "GET", url, false ); // false for synchronous request
            xmlHttp.send();
            console.log(xmlHttp.responseText);
        }

        //setTimeout(update(deltaTime), 600000);

        myVideo.addEventListener('timeupdate', function() {
            console.log(this.currentTime);
            console.log(deltaTime);

            if (this.currentTime < start) {
                this.currentTime = start;
            }

            deltaTime += this.currentTime - lastTime;
            lastTime = this.currentTime;

            if (this.currentTime > end) {
                this.currentTime = start;
                myVideo.play();
                update(end-start);
            }
        }, false);

        

        myVideo.addEventListener('ended', function() {
            this.currentTime = start;
            myVideo.play();
            update(this.duration-start);
        },false);

        myVideo.play();
    </script>

    <br>
    <?php include "history.php" ?>
</body>

</html>