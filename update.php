<?php

use Chassit\Repeat\DB;

require_once "init.php";

if ($_SERVER['REQUEST_METHOD'] == "POST") {
    $input = json_decode(file_get_contents("php://input"));
    $name = $input->v;
    $time = $input->t;
    $start = $input->s ?? null;
    $end = $input->e ?? null;

    if ($start !== null)
        $start = is_numeric($start) ? floatval($start) : null;
    if ($end !== null)
        $end = is_numeric($end) ? floatval($end) : null;

    if ($end > 90000){ exit("Playtime error."); }
    //if ($name == "OKWVNeDYZmU"){ exit("Staven Misshandlar"); }

    $file = glob('files/*-' . $name . '.mp4');

    if (count($file) == 1) {
        $insertOneResult = DB::getRepeatCollection()->updateOne(
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

} else {
    http_response_code(400);
    echo "Unsupported";
}
