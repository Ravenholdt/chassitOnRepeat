<?php
    ini_set("display_errors", 1);
    error_reporting(E_ALL);

    use Dotenv\Dotenv;
    use MongoDB\Driver\ServerApi;

    require __DIR__ . '/vendor/autoload.php';

    $dotenv = Dotenv::createImmutable(__DIR__);
    $dotenv->load();

    $serverApi = new ServerApi(ServerApi::V1);
    $client = new MongoDB\Client(
        'mongodb+srv://'.$_ENV["MONGO_USER"].':'.$_ENV["MONGO_PASS"].'@null.t2drt9o.mongodb.net/?retryWrites=true&w=majority', [], ['serverApi' => $serverApi]);

    $collection = $client->repeat->data;

    $filter  = [];
    $options = ['sort' => ['playtime' => -1]];
    $cursor = $collection->find($filter, $options);

    foreach ( $cursor as $id => $value )
    {
        $start = "";
        $end = "";
        $repeat = "";

        if (isset($value['start'])){ $start = "&s=" . $value['start']; }
        if (isset($value['end'])){ $end = "&e=" . $value['end']; }
        if (isset($value['playtime'])){ $repeat = $value['playtime']; }

        $file = glob('files/*-'. $value['name'] .'.mp4');
        echo "<a href=?v=" . $value['name'] . $start . $end . 
        ">";
        echo $repeat . " : ";
        echo explode("/", $file[0])[1];
        echo "</a>";
        echo "<br>";
    }

    //var_dump($history);
?>