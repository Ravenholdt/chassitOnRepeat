<?php
    ini_set("display_errors", 1);
    error_reporting(E_ALL);

    include "dbconnect.php";

    $name = $_GET["v"];
    $time = $_GET["t"];
    $start = $_GET["s"];
    $end = $_GET["e"];

    $file = glob('files/*-'. $name .'.mp4');

    if (count($file) == 1){
        $collection = $client->repeat->data;
    
        $insertOneResult = $collection->updateOne(
            ['name' => $name],
            [
            '$set' => [
                'lastplayed' => time(),
                'start' => $start,
                'end' => $end
            ], 
            '$inc' => [
                'playtime' => (int)$time]
            ],
            ['upsert' => true]
        );
    }

    //$deleteResult = $collection->deleteMany([]);
?>