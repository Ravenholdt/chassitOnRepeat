<?php
    ini_set("display_errors", 1);
    error_reporting(E_ALL);
    
    use Dotenv\Dotenv;
    use MongoDB\Driver\ServerApi;

    require __DIR__ . '/vendor/autoload.php';

    $dotenv = Dotenv::createImmutable(__DIR__);
    $dotenv->load();

    $name = $_GET["v"];
    $time = $_GET["t"];

    $file = glob('files/*-'. $name .'.mp4');
    
    if (count($file) == 1){
        $serverApi = new ServerApi(ServerApi::V1);
        $client = new MongoDB\Client(
            'mongodb+srv://'.$_ENV["MONGO_USER"].':'.$_ENV["MONGO_PASS"].'@null.t2drt9o.mongodb.net/?retryWrites=true&w=majority', [], ['serverApi' => $serverApi]);
    
        $collection = $client->repeat->data;
    
        $insertOneResult = $collection->updateOne(
            ['name' => $name],
            ['$set' => ['lastplayed' => time()], '$inc' => ['playtime' => (int)$time]],
            ['upsert' => true]
        );
    }

    //$deleteResult = $collection->deleteMany([]);
?>