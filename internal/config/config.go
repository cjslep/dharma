// dharma is a supplementary corporation community tool for Eve Online.
// Copyright (C) 2021 Cory Slep
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package config

type Config struct {
	ClientID             string `ini:"dharma_client_id" comment:"The client identifier CCP Games gives your application when registering on the ESI site, to identify your particular software instance."`
	APIKey               string `ini:"dharma_api_key" comment:"The secret CCP Games gives your application to verify the authenticity of your software instance."`
	EnableConsoleLogging bool   `ini:"dharma_debug_console_log" comment:"When true, logs directly to console, which is best used during software development."`
	LogDir               string `ini:"dharma_log_directory" comment:"Directory location to write log files to, which can be useful when filing bug reports. (default: ./)"`
	LogFile              string `ini:"dharma_log_file" comment:"Name of the log files, which can rotate and use this as a base naming convention. (default: dharma.log)"`
	NLogFiles            int    `ini:"dharma_n_log_files_rotation" comment:"Number of log files to keep on hand at a time, before rotating and overwriting them, zero means don't rotate. (default: 5)"`
	MaxMBSizeLogFiles    int    `ini:"dharma_max_mb_size_log_files_rotation" comment:"Max size in megabytes of a single log file, before rotating to a new file. (default: 100)"`
	MaxDayAgeLogFiles    int    `ini:"dharma_max_age_days_log_files_rotation" comment:"Max number of days to keep a single log file, before rotating it, zero means don't rotate based on age. (default: 0)"`

	NPreview     int `ini:"dharma_length_post_preview" comment:"Number of preview texts to display per tag (default: 3)"`
	LenPreview   int `ini:"dharma_length_post_preview" comment:"The length of preview text to display (default: 80)"`
	MaxHTMLDepth int `ini:"dharma_max_html_parsing_depth" comment:"The deepest HTML parsing allowed before abandoning (default: 255)"`
}
