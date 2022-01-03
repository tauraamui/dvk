package stdlibx

var Modules = map[string]string{
	"metadata": `
text := import("text")

regexp_within_braces := ;;;;\[[^]]*\];;;; // matches all characters between '[' ']'

export func(log_line) {
    meta_data := { time: undefined, level: undefined, full_log: log_line, log_text: undefined }

    matches := text.re_find(regexp_within_braces, log_line, 2)
    if matches != undefined {
        first_match_start_index := undefined
        first_match_end_index := undefined
        if len(matches) >= 1 {
            if m := matches[0]; len(m) >= 1 {
                m = m[0]
                first_match_start_index = m.begin
                first_match_end_index = m.end
                meta_data.time = text.trim_suffix(text.trim_prefix(m.text, "["), "]")
            }
        }

        second_match_start_index := undefined
        second_match_end_index := undefined
        if len(matches) >= 2 {
            if m := matches[1]; len(m) >= 1 {
                m = m[0]
                second_match_start_index = m.begin
                second_match_end_index = m.end
                meta_data.level = text.trim_suffix(text.trim_prefix(m.text, "["), "]")
            }
        }


        if first_match_start_index != undefined && first_match_end_index != undefined {
            // TODO(tauraamui): Re-think log content selecting
            // this is too destructive and still assumes that the matches
            // proceed the remaining log line's content, which is not certain
            log_line = log_line[first_match_start_index:][first_match_end_index:]
            if second_match_start_index != undefined && second_match_end_index != undefined {
                second_match_start_index = second_match_start_index - first_match_end_index
                second_match_end_index = second_match_end_index - first_match_end_index
            }
        }

        if second_match_start_index != undefined && second_match_end_index != undefined {
            log_line = log_line[second_match_start_index:][second_match_end_index:]
        }

        meta_data.log_text = log_line
    }

    return meta_data
}

	`,
}
