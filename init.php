<?php

use Dotenv\Dotenv;

require __DIR__ . '/vendor/autoload.php';

$dotenv = Dotenv::createImmutable(__DIR__);
$dotenv->load();

if ($_ENV["DEBUG"] ?? false)
{
    ini_set("display_errors", 1);
    error_reporting(E_ALL);
}

