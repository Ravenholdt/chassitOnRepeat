<html>

<?php
    $name = $_GET["v"];

    $file = glob('files/*-'. $name .'.mp4');

    $start = $_GET["s"];
    if (!isset($_GET["s"])){ $start = 0; }

    $end = $_GET["e"];
    //if (!isset($_GET["s"])){ $end = 0; }
?>

<header>
</header>

<body>
    <video id="myvideo" controls autoplay>
        <source src="<?php echo $file[0]; ?>" type="video/mp4">
        Your browser does not support HTML5 video.
    </video>

    <script>
        myVideo = document.getElementById("myvideo")
        /*
        document.getElementById('myvideo').addEventListener('loadedmetadata', function() {
            this.currentTime = start; 
        }, false);
        */

        function update(t) {
            var xmlHttp = new XMLHttpRequest();

            var url="update.php?v=<?php echo $file[0]; ?>&t="+t;
            xmlHttp.open( "GET", url, false ); // false for synchronous request
            xmlHttp.send();
            console.log(xmlHttp.responseText);
        }

        myVideo.addEventListener('timeupdate', function() {
            console.log(this.currentTime);
            if (this.currentTime < <?php echo $start ?>) {
                this.currentTime = <?php echo $start ?>;
            }
            if (this.currentTime > <?php echo $end ?>) {
                this.currentTime = <?php echo $start ?>;
                update(<?php echo $end-$start ?>);
            }
        }, false);
        document.getElementById("myvideo").play()
    </script>
</body>

</html>