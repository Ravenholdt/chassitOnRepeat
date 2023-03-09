<?php

namespace Chassit\Repeat;

use MongoDB\Driver\Exception\ConnectionTimeoutException;

class History
{
    /**
     * @var Video[]
     */
    static ?array $allVideos = null;

    static bool $noCloud = false;

    /**
     * Returns a list of videos found on the files directory
     * @return Video[]
     */
    private static function getFiles(): array
    {
        $videos = [];
        foreach (glob('files/*.mp4') as $filename) {

            // Matches the name and the if of the file
            // The if is the last 11 characters before ".mp4"
            if (preg_match("/^files\/(.*)-([A-Za-z0-9_-]{11}).mp4$/", $filename, $matches))
                $videos[$matches[2]] = new Video($matches[2], $matches[1]);
        }

        return $videos;
    }

    public static function getRandomVideo(): Video
    {
        $files = array_values(self::getRepeats(self::getFiles()));
        $video = $files[rand(0, sizeof($files))];

        /*
        // djNull
        foreach ($files as $key => $value) {
            //echo $key . " " . $value->name . " " . $value->id . " " . $value->playtime . "<br>";
            if ($value->id == "JE37e1eK2mY"){
                return $value;
            }
        }
        */
        
        
        return $video;
    }

    /**
     * @param array|object $arr
     * @param string $key
     * @return float|null
     */
    private static function getFloatVal($arr, string $key): ?float
    {
        if (isset($arr[$key]))
        {
            $val = $arr[$key];
            return is_numeric($val) ? floatval($val): null;
        }

        return null;
    }

    /**
     * Returns a list of videos with updated repeats data
     * @param Video[] $videos A list of videos to merge the database data with
     * @return Video[]
     */
    public static function getRepeats(array $videos): array
    {
        $cursor = DB::getRepeatCollection()->find([], ['sort' => ['playtime' => -1]]);
        foreach ($cursor as $value) {
            $id = $value["name"];

            // Check if video is available, update data if it is otherwise add a placeholder...
            if (isset($videos[$id])) {
                $video = $videos[$id];
                $video->start = self::getFloatVal($value, 'start');
                $video->end = self::getFloatVal($value, 'end');
                $video->playtime = $value["playtime"] ?? 0;

                $videos[$id] = $video;
            } else { // ...but it's without the placeholder.
                //$videos[$id] = new Video($value["name"], $value["name"], self::getFloatVal($value, 'start'), self::getFloatVal($value, 'end'), $value["playtime"] ?? 0);
            }
        }
        return $videos;
    }

    /**
     * Loads all the videos from file and the database
     * @return void
     */
    private static function loadVideos()
    {
        $files = self::getFiles();
        try {
            self::$allVideos = self::getRepeats($files);
        } catch (ConnectionTimeoutException $exception) {
            // If connection has failed just use the normal files without time and interval
            self::$allVideos = $files;
            self::$noCloud = true;
        }
    }

    public static function getTotalTime(): int
    {
        if (self::$allVideos == null)
            self::loadVideos();

        $time = 0;
        foreach (self::$allVideos as $video) $time += $video->playtime;
        return $time;
    }

    public static function toDisplayTime(int $playtime, bool $showDays = false): string
    {
        $secondsInAMinute = 60;
        $secondsInAnHour  = 3600;
        $secondsInADay    = 86400;

        $hourSeconds = $showDays ? $playtime % $secondsInADay: $playtime;
        $minuteSeconds = $hourSeconds % $secondsInAnHour;
        $remainingSeconds = $minuteSeconds % $secondsInAMinute;

        $days = floor($playtime / $secondsInADay);
        $hours = floor($hourSeconds / 3600);
        $minutes = floor($minuteSeconds / $secondsInAMinute);
        $seconds = ceil($remainingSeconds);

        return $days > 0 && $showDays ?
            sprintf("%dd %02dh %02dm %02ds", $days, $hours, $minutes, $seconds):
            ($hours > 0 ?
                sprintf("%dh %02dm %02ds", $hours, $minutes, $seconds) :
                ($minutes > 0 ?
                    sprintf("%dm %02ds", $minutes, $seconds) :
                    sprintf("%ds", $seconds)));
    }

    /**
     * Renders a list of links to videos
     * @return void
     */
    public static function render()
    {
        if (self::$allVideos == null)
            self::loadVideos();

        // Sort videos before displaying
        usort(self::$allVideos, function (Video $a, Video $b) {
            return $a->playtime == 0 && $b->playtime == 0 ?
                strtolower($a->name) <=> strtolower($b->name) : // If playtime is zero in both compare on name
                $b->playtime <=> $a->playtime; // Otherwise compare on playtime
        });

        foreach (self::$allVideos as $video) {
            $args = "";
            /*if ($video->start !== null && $video->start > 0) {
                $args .= "&s=$video->start";
            }
            if ($video->end !== null && $video->end > 0) {
                $args .= "&e=$video->end";
            }*/

            $time = $video->playtime > 0 ? self::toDisplayTime($video->playtime): "";

            // If Yoda torture has been playing for more than
            // 200 hours, play other yoda instead.
            if ($video->id == "OKWVNeDYZmU" && $video->playtime > 720000) {
                $video->id = "IUMCyAR6U0";
                $args = "&s=13.4&e=243";
            }

            echo "<a href=\"?v=$video->id$args\"> $time : $video->name</a><br>";
        }
    }
}
