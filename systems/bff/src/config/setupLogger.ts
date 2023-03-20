import { format, transports, createLogger, Logger } from "winston";
const { combine, timestamp, errors, label, printf, splat, colorize } = format;

// Handles pretty printing objects passed to logger rather than `[object Object]
// Does not handle string + object log message. To log this way, use %o, such as logger.info("some message %o", objectToPrint)
// Ref: https://github.com/winstonjs/winston/issues/1217

const customFormat = printf(info => {
    if (typeof info.message === "object") {
        info.message = JSON.stringify(info.message, null, 2);
    }
    let message = `${info.timestamp} ${info.level} [${info.label}]: ${info.message}`;
    if (info.meta) {
        message += ` ${info.meta.res?.statusCode || ""}`;
    }

    return message;
});

const options = {
    combinedFile: {
        level: "info",
        filename: "logs/app.log",
        handleExceptions: true,
        json: true,
        maxsize: 5242880, // 5MB
        maxFiles: 5,
        colorize: false,
    },
    erorrFile: {
        level: "error",
        filename: "logs/errors.log",
        handleExceptions: true,
        json: true,
        maxsize: 5242880, // 5MB
        maxFiles: 5,
        colorize: false,
    },
    console: {
        level: "debug",
        handleExceptions: true,
        json: false,
        colorize: true,
    },
    exceptions: {
        filename: "logs/exceptions.log",
    },
};

// Accepts label to be prepended for logs produced by logger and returns logger
// Example:
// const logger = require('../config/logger')('userController');
// logger.info('some log')

const setupLogger = (sourceLabel: string): Logger => {
    return createLogger({
        transports: [
            new transports.File(options.combinedFile),
            new transports.File(options.erorrFile),
            new transports.Console({
                format: combine(
                    label({ label: sourceLabel }),
                    colorize(),
                    timestamp(),
                    errors({ stack: true }),
                    splat(),
                    customFormat,
                ),
            }),
        ],
        format: combine(
            label({ label: sourceLabel }),
            colorize(),
            timestamp(),
            errors({ stack: true }),
            splat(),
            customFormat,
        ),
        exitOnError: false,
        exceptionHandlers: [
            new transports.File(options.exceptions),
            new transports.Console({
                format: combine(
                    label({ label: sourceLabel }),
                    colorize(),
                    timestamp(),
                    errors({ stack: true }),
                    splat(),
                    customFormat,
                ),
            }),
        ],
    });
};

export default setupLogger;
