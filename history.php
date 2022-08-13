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
        $DBstart = "";
        $DBend = "";
        $DBrepeat = "";

        if (isset($value['start'])){ $DBstart = "&s=" . $value['start']; }
        if (isset($value['end'])){ $DBend = "&e=" . $value['end']; }
        if (isset($value['playtime'])){ $DBrepeat = $value['playtime']; }

        $file = glob('files/*-'. $value['name'] .'.mp4');
        echo "<a href=?v=" . $value['name'] . $DBstart . $DBend . 
        ">";
        echo $DBrepeat . " : ";
        echo explode(".mp4", explode("/", $file[0])[1])[0];
        echo "</a>";
        echo "<br>\n";
    }

    //var_dump($history);
?>