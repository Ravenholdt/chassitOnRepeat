/**
 * Format a number to the given amount of digits
 * @param {number} num
 * @param {number} digits
 */
function formatNumber(num, digits = 2) {
    return num.toLocaleString('en-US', {minimumIntegerDigits: digits, useGrouping: false})
}

/**
 * Formats the time to show days, hours, minutes and seconds
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