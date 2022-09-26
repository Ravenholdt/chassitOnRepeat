<?php

namespace Chassit\Repeat;

class Video
{
    public string $name;
    public string $id;
    public ?string $start;
    public ?string $end;
    public int $playtime;

    public function __construct(string $id, string $name, string $start = null, string $end = null, int $playtime = 0)
    {
        $this->name = $name;
        $this->id = $id;
        $this->start = $start;
        $this->end = $end;
        $this->playtime = $playtime;
    }
}