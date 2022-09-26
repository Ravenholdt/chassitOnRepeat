<?php

namespace Chassit\Repeat;

class History
{
    /**
     * @var Video[]
     */
    static array $allVideos = [];

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
            preg_match("/^files\/(.*)-([A-Za-z0-9_-]{11}).mp4$/", $filename, $matches);
            $videos[$matches[2]] = new Video($matches[2], $matches[1]);
        }

        return $videos;
    }

    /**
     * Returns a list of videos with updated repeats data
     * @param Video[] $videos A list of videos to merge the database data with
     * @return Video[]
     */
    private static function getRepeats(array $videos): array
    {
        $cursor = DB::getRepeatCollection()->find([], ['sort' => ['playtime' => -1]]);
        foreach ($cursor as $value) {
            $id = $value["name"];

            // Check if video is available, update data if it is otherwise add a placeholder
            if (isset($videos[$id])) {
                $video = $videos[$id];
                $video->start = $value["start"] ?? null;
                $video->end = $value["end"] ?? null;
                $video->playtime = $value["playtime"] ?? 0;

                $videos[$id] = $video;
            } else {
                $videos[$id] = new Video($value["name"], $value["name"], $value["start"] ?? null, $value["end"] ?? null, $value["playtime"] ?? 0);
            }
        }
        return $videos;
    }

    /**
     * Loads all the videos from file and the database
     * @return void
     */
    public static function loadVideos()
    {
        $files = self::getFiles();
        self::$allVideos = self::getRepeats($files);
    }

    public static function getTotalTime(): int
    {
        $time = 0;
        foreach (self::$allVideos as $video) $time += $video->playtime;
        return $time;
    }

    /**
     * Renders a list of links to videos
     * @return void
     */
    public static function render()
    {
        // Sort videos before displaying
        usort(self::$allVideos, function (Video $a, Video $b) {
            return $a->playtime == 0 && $b->playtime == 0 ?
                strtolower($a->name) <=> strtolower($b->name) : // If playtime is zero in both compare on name
                $b->playtime <=> $a->playtime; // Otherwise compare on playtime
        });

        foreach (self::$allVideos as $video) {
            $args = "";
            if ($video->start !== null) {
                $args .= "&s=$video->start";
            }
            if ($video->end !== null) {
                $args .= "&e=$video->end";
            }

            $time = $video->playtime > 0 ? $video->playtime : "";
            echo "<a href=\"?v=$video->id$args\"> $time : $video->name</a><br>";
        }
    }
}