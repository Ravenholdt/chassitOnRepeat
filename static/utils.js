/**
 * Shorter function to get the element by id
 * @param {string} id The id to locate with
 * @returns {HTMLElement} The corresponding element
 */
function elementById(id) {
    return document.getElementById(id)
}

/**
 * Shorter function to get the elements from a class
 * @param {string} classNames The string of the classnames to locate
 * @returns {HTMLCollectionOf<Element>}
 */
function elementsByClass(classNames) {
    return document.getElementsByClassName(classNames);
}

/**
 * Format a number to the given number of digits
 * @param {number} num
 * @param {number} digits
 */
function formatNumber(num, digits = 2) {
    return num.toLocaleString('en-US', {minimumIntegerDigits: digits, useGrouping: false})
}

/**
 * Formats the time to show days, hours, minutes, and seconds
 * @param {number} time
 */
function formatTime(time) {
    const hourSeconds = time % (60 * 60 * 24);
    const minuteSeconds = hourSeconds % (60 * 60)
    const remainingSeconds = minuteSeconds % 60

    const days = Math.floor(time / (60 * 60 * 24))
    const hours = Math.floor(hourSeconds / (60 * 60))
    const minutes = Math.floor(minuteSeconds / 60)
    const seconds = Math.floor(remainingSeconds)

    if (days > 0)
        return `${days}d ${formatNumber(hours)}h ${formatNumber(minutes)}m ${formatNumber(seconds)}s`;
    if (hours > 0)
        return `${formatNumber(hours)}h ${formatNumber(minutes)}m ${formatNumber(seconds)}s`
    if (minutes > 0)
        return `${formatNumber(minutes)}m ${formatNumber(seconds)}s`
    return `${seconds}s`
}

/**
 * Returns a function that checks for it the video has ended
 * @param {HTMLVideoElement} video The video element
 * @param {number} start The start time of the video
 * @param {number} end The end time of the video
 * @param {(function(number): Promise<void>)} endHandler A callback function to be called when video has reached the end
 * @returns {(function(): Promise<void>)|*}
 */
function updateCheckFunc(video, start, end, endHandler) {
    return async function updateTimer() {
        if (video.currentTime < start) {
            video.currentTime = start;
            console.log("Jump Forward");
        }

        if (video.currentTime >= end) {
            console.log("Restart");
            await endHandler(end - start);
        }

        if (video.currentTime >= video.duration) {
            await endHandler(video.duration - start);
        }
    }
}

/**
 * Creates a UUID v4
 *
 * Copied from https://stackoverflow.com/a/2117523/2230267
 * @returns {string}
 */
function createUUIDv4() {
    return "10000000-1000-4000-8000-100000000000".replace(/[018]/g,
            c =>
                (+c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> +c / 4)
                .toString(16)
    );
}