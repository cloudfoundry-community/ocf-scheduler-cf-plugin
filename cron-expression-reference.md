# Cron Expression Reference

The OCF Scheduler uses `cron expression` notation for all schedules. Expressions are parsed by the [go-cron](https://github.com/netresearch/go-cron) library using its standard parser, which supports standard cron syntax, extended day-of-month/day-of-week syntax, descriptor shortcuts, and timezone prefixes.

## Standard Cron Fields

The core format are 5 required fields:

```
minute hour day-of-month(DOM) month day-of-week(DOW)
```

| Field | Required | Range | Named Values |
|---|---|---|---|
| Minute | Yes | `0-59` | None |
| Hour | Yes | `0-23` | None |
| Day of Month | Yes | `1-31` | None |
| Month | Yes | `1-12` | `JAN`-`DEC` (case-insensitive) |
| Day of Week | Yes | `0-7` (0 and 7 both represent Sunday) | `SUN`-`SAT` (case-insensitive) |

### Operators

Operators can be used within any field to define schedules beyond a single fixed value.

| Syntax | Meaning | Example |
|---|---|---|
| `*` | Every value in the field's range | `* * * * *` — every minute |
| `N` | A specific value | `30 * * * *` — at minute 30 |
| `N-M` | An inclusive range; wraps around when N > M | `1-5` in day-of-week — Monday through Friday; `22-2` in hours — 22, 23, 0, 1, 2 |
| `N-M/S` | A range with a step | `0-30/5` in minutes — every 5 minutes from 0 to 30 |
| `*/S` | Wildcard with a step | `*/15` — every 15 units |
| `N/S` | Start at N, step by S to the field maximum | `5/10` in minutes — 5, 15, 25, 35, 45, 55 |
| `N,M,...` | A list of values | `1,15` in day-of-month — 1st and 15th |

### Examples

```
*/5 * ? * *            Every 5 minutes
30 14 * * *            2:30 PM daily
0 0 1 * *              Midnight on the 1st of every month
0 9 * * MON-FRI        9:00 AM on weekdays
30 */2 * * *           Minute 30 of every 2nd hour
0 0 1,15 * *           Midnight on the 1st and 15th of every month
0 8 * * MON,WED,FRI    8:00 AM on Monday, Wednesday, and Friday
0 22-6 * * *           Every hour from 10 PM to 6 AM (overnight)
0 0 * OCT-MAR *        Midnight daily during October through March
0-30/10 9 * * *        Minutes 0, 10, 20, 30 of 9 AM
0 9 L * *              9:00 AM on the last day of every month
0 9 15W * *            9:00 AM on the nearest weekday to the 15th
0 17 LW * *            5:00 PM on the last weekday of every month
0 10 * * FRI#1         10:00 AM on the first Friday of every month
0 12 * * THU#L         Noon on the last Thursday of every month
```

## Extended Day-of-Month Syntax

Extended operators are available in the day-of-month field for expressing relative and weekday-aware dates.

| Syntax | Meaning | Example |
|---|---|---|
| `?` | Same as `*` (wildcard) | `0 0 ? * MON` — midnight every Monday |
| `L` | Last day of the month | `0 0 L * *` — midnight on the last day |
| `L-N` | Nth day from the end of the month (N: 1-30) | `0 0 L-3 * *` — 3rd to last day |
| `NW` | Nearest weekday to the Nth day | `0 0 15W * *` — nearest weekday to the 15th |
| `LW` | Last weekday of the month | `0 0 LW * *` — last weekday |

These can be combined with commas and standard values: `1,15,L` means the 1st, 15th, and last day of the month.

## Extended Day-of-Week Syntax

Extended operators are available in the day-of-week field for expressing specific occurrences of a weekday within a month.

| Syntax | Meaning | Example |
|---|---|---|
| `?` | Same as `*` (wildcard) | `0 0 1 * ?` — midnight on the 1st of every month |
| `DOW#N` | Nth occurrence of a weekday in the month (N: 1-5) | `FRI#3` — 3rd Friday |
| `DOW#L` | Last occurrence of a weekday in the month | `FRI#L` — last Friday |

These can be combined with commas and standard values: `MON,FRI#3,SUN#L` means every Monday, the 3rd Friday, and the last Sunday.

## Descriptor Shortcuts

Predefined schedules can be specified with an `@` prefix instead of writing out cron fields.

| Descriptor | Equivalent | Meaning |
|---|---|---|
| `@yearly` | `0 0 1 1 *` | Midnight, January 1st |
| `@annually` | `0 0 1 1 *` | Same as `@yearly` |
| `@monthly` | `0 0 1 * *` | Midnight, 1st of every month |
| `@weekly` | `0 0 * * 0` | Midnight every Sunday |
| `@daily` | `0 0 * * *` | Midnight every day |
| `@midnight` | `0 0 * * *` | Same as `@daily` |
| `@hourly` | `0 * * * *` | Start of every hour |
| `@every <duration>` | N/A | Fixed interval (e.g. `@every 1h30m`, `@every 5m`, `@every 30s`) |

The `@every` descriptor accepts any Go duration string: `s` (seconds), `m` (minutes), `h` (hours), and combinations such as `1h30m`.

## Timezone Support

By default, schedules are evaluated in the scheduler server's local timezone. There are two ways to specify a different timezone.

### Using the `--timezone` flag

The `schedule-job` and `schedule-call` commands accept a `--timezone` (`-t`) flag:

```
cf schedule-job my-job --timezone America/New_York "0 9 * * *"
cf schedule-call my-call -t Europe/Berlin "30 14 * * MON-FRI"
```

This is the recommended approach as it keeps the cron expression and timezone clearly separated.

### Prefixing the cron expression with an environment variable

Alternatively, the cron expression can be prefixed with a `TZ=` or `CRON_TZ=` environment variable, the same way a shell command can be prefixed with a variable assignment:

```
CRON_TZ=America/New_York 0 9 * * *
TZ=Europe/Berlin 30 14 * * MON-FRI
```

Both `CRON_TZ=` and `TZ=` are supported and behave identically.

> **Do not combine both methods.** If you use the `--timezone` flag, do not also include a `TZ=` or `CRON_TZ=` prefix in the expression. The flag works by prepending `CRON_TZ=` to the expression, so using both would produce a malformed expression with two timezone prefixes.

### Listing available timezones

Run `cf scheduler-time-zones` (alias: `stz`) to retrieve the list of supported timezones from the scheduler. The output includes each timezone's name, whether it observes daylight saving time (DST), any aliases, and which timezone the server itself is using.

The timezone value must be a valid [IANA timezone name](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).

### Daylight saving time

The scheduler handles DST transitions automatically:

- **Spring forward** (clocks skip an hour) — jobs scheduled during the skipped hour run immediately after the transition.
- **Fall back** (clocks repeat an hour) — jobs scheduled during the repeated hour run once, during the first occurrence.

To avoid surprises around DST transitions, schedule time-sensitive jobs outside the typical 1:00–3:00 AM transition window, or use a timezone that does not observe DST (e.g. `UTC`, `America/Phoenix`).
