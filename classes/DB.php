<?php

namespace Chassit\Repeat;

use MongoDB\Client;
use MongoDB\Collection;
use MongoDB\Driver\ServerApi;

class DB
{
    private static $_client;

    /**
     * Creates a new instance of the MongoDB client
     * @return void
     */
    private static function setup()
    {
        $serverApi = new ServerApi(ServerApi::V1);
        self::$_client = new Client(
            'mongodb+srv://' . $_ENV["MONGO_USER"] . ':' . $_ENV["MONGO_PASS"] . '@null.t2drt9o.mongodb.net/?retryWrites=true&w=majority', [], ['serverApi' => $serverApi]);
    }

    /**
     * Returns the MongoDB client
     * @return Client
     */
    public static function client(): Client
    {
        if (self::$_client == null)
            self::setup();
        return self::$_client;
    }

    /**
     * Helper function to get the repeat collection
     * @return Collection
     */
    public static function getRepeatCollection(): Collection {
        return DB::client()->selectCollection("repeat", "data");
    }
}