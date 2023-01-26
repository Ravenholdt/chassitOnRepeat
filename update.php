<?php

use Chassit\Repeat\DB;

require_once "init.php";

if ($_SERVER['REQUEST_METHOD'] == "POST") {
    $input = json_decode(file_get_contents("php://input"));
    $name = $input->v;
    $time = $input->t ?? 0;
    $start = $input->s ?? null;
    $end = $input->e ?? null;

    if ($time > 90000){ http_response_code(400); exit("Playtime error."); }
    //if ($name == "OKWVNeDYZmU"){ http_response_code(418); exit("Staven Misshandlar"); }

    $update = [];

    if ($time !== 0)
        $update['lastplayed'] = time();
    if ($start !== null)
        $update['start'] = is_numeric($start) ? floatval($start) : null;
    if ($end !== null)
        $update['end'] = is_numeric($end) ? floatval($end) : null;

    if (array_key_exists('end', $update) && array_key_exists('start', $update)) {
        if ($update['end'] <= $update['start']) {
            http_response_code(400);
            exit("Interval error.");
        }
    }

    $file = glob('files/*-' . $name . '.mp4');

    if (count($file) == 1) {
        $insertOneResult = DB::getRepeatCollection()->updateOne(
            ['name' => $name],
            [
                '$set' => $update,
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
