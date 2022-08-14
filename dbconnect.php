<?php
    use Dotenv\Dotenv;
    use MongoDB\Driver\ServerApi;

    require __DIR__ . '/vendor/autoload.php';

    $dotenv = Dotenv::createImmutable(__DIR__);
    $dotenv->load();

    $serverApi = new ServerApi(ServerApi::V1);
    $client = new MongoDB\Client(
        'mongodb+srv://'.$_ENV["MONGO_USER"].':'.$_ENV["MONGO_PASS"].'@null.t2drt9o.mongodb.net/?retryWrites=true&w=majority', [], ['serverApi' => $serverApi]);
?>