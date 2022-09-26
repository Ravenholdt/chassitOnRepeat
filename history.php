<?php
    ini_set("display_errors", 1);
    error_reporting(E_ALL);

    $collection = $client->repeat->data;

    $totalTime = 0;

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
        if (isset($value['playtime'])){ $DBrepeat = $value['playtime']; $totalTime += $DBrepeat; }

        $file = glob('files/*-'. $value['name'] .'.mp4');
        echo "<a href=?v=" . $value['name'] . $DBstart . $DBend . ">";
        echo $DBrepeat . " : ";
        echo explode(".mp4", explode("/", $file[0])[1])[0];
        echo "</a>";
        echo "<br>\n";
    }

    //var_dump($history);
?>