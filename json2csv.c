/*
json2csv converts json data to csv format

http://github.com/jehiah/json2csv

copyright 2010 Jehiah Czebotar <jehiah@gmail.com> 

Version: 1.0

uses json-c fork at http://github.com/jehiah/json-c
*/

#include <sys/types.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "json/json.h"

#define SUCCESS 0
#define FAILURE 1

#define JSON_GET_STR(json_obj, field) (json_object_to_json_string(json_object_object_get(json_obj, field)))
#define JSON_GET_INT(json_obj, field) (json_object_get_int(json_object_object_get(json_obj, field)))
#define JSON_FREE(json_obj) (json_object_put(json_obj))
#define JSON_DEBUG 0

#define BUFF_SIZE 20480

int parse_fields(char *str, char **field_array);
void parse_output_keys(char *str);
void process_line(char *source, FILE *output);
void run(char *input_filename, char *output_filename);
void usage();

char *json_string_to_csv_string(const char *json_string);

int parse_json(char *json_str, struct json_object **json_obj); 


int verbose=0;
/* variables to hold the keys we want to output */
static char *output_keys[64];
static int  num_output_keys = 0;
char buffer[BUFF_SIZE];


int parse_json(char *json_str, struct json_object **json_obj)
{
    if (!json_str) {
        if(JSON_DEBUG) fprintf(stderr, "ERR: empty message\n");
        return FAILURE;
    }

    if (strlen(json_str) < 3) {
        if(JSON_DEBUG) fprintf(stderr, "ERR: empty json\n");
        return FAILURE;
    }

    *json_obj = json_tokener_parse(json_str);

    if (*json_obj == NULL) {
        fprintf(stderr, "ERR: unable to parse json (%s) \n", json_str);
        return FAILURE;
    }

    return SUCCESS;
}

void process_line(char *source, FILE *output)
{
    int i=0;
    struct json_object *json_obj;
    if (parse_json(source, &json_obj) == FAILURE) {
        return;
    }
    for (i=0; i < num_output_keys; i++) {
        char *str = json_string_to_csv_string(JSON_GET_STR(json_obj, output_keys[i]));
        if (i == 0) {
            fprintf(output, "%s", strcmp(str, "null") == 0 ? "" : str);
        } else {
            fprintf(output, ",%s", strcmp(str, "null") == 0 ? "" : str);
        }
        if (str) {
            free(str);
        }
    }
    fprintf(output,"\n");
    if (json_obj) {
        JSON_FREE(json_obj);
    }
}

void run(char *input_filename, char *output_filename) {
    FILE *in_file = stdin;
    FILE *out_file = stdout;
    struct json_object *data;
    const char *json;
    ssize_t read;
    size_t line_len=0;

    if (input_filename) {
        if (verbose) fprintf(stderr, "input file is %s\n", input_filename);
        in_file = fopen(input_filename, "r" );
    } else {
        if (verbose) fprintf(stderr, "input is stdin\n");
    }
    if (output_filename) {
        if (verbose) fprintf(stderr, "output file is %s\n", output_filename);
        out_file = fopen(output_filename, "w" );
    } else {
        if (verbose) fprintf(stderr, "output is stdout\n");
    }
    
    int lines = 0;
    if (!in_file){
        perror(input_filename);
        return;
    }
    if (!out_file){
        perror(output_filename);
        return;
    }
    data = json_object_new_object();
    while (fgets(buffer, BUFF_SIZE, in_file)) {
        lines +=1;
        process_line(buffer, out_file);
    }
    if (verbose) fprintf(stderr, "processed %d lines\n", lines);
    fclose(in_file);
    fclose(out_file);
    JSON_FREE(data);

    int i;
    for (i = 0; i < num_output_keys; i++) {
        free(output_keys[i]);
    }
}


void parse_output_keys(char *str)
{
    int i;
    num_output_keys = parse_fields(str, output_keys);

    if (verbose) {
        for (i=0; i < num_output_keys; i++) {
            fprintf(stderr, "selecting key: \"%s\"\n", output_keys[i]);
        }
    }
    return;
}

/*
 * Parse a comma-delimited list of strings and put them
 * in an char array. Array better have enough slots
 * because I didn't have time to work out the memory allocation.
 */
int parse_fields(char *str, char **field_array)
{
    int i;
    const char delim[] = ",";
    char *tok, *str_ptr, *save_ptr;

    if (!str) return;

    str_ptr = strdup(str);

    tok = strtok_r(str_ptr, delim, &save_ptr);

    i = 0;
    while (tok != NULL) {
        field_array[i] = strdup(tok);
        tok = strtok_r(NULL, delim, &save_ptr);
        i++;
    }
    if (str_ptr) {
        free(str_ptr);
    }
    return i;
}


/*
 * Transforms the JSON escape sequence \" into the CSV escape sequence "".
 */
char *json_string_to_csv_string(const char *json_string)
{
    char *str = strdup(json_string);
    char *head = str;
    while (*str) {
        if (*str == '\\' && *(str+1) == '\"') {
            *str = '\"';
        }
        ++str;
    }
    return head;
}

void usage(){
    fprintf(stderr, "usage: json2csv\n");
    fprintf(stderr, "\t-k fields,to,output\n");
    fprintf(stderr, "\t-i /path/to/input.json (optional; default is stdin)\n");
    fprintf(stderr, "\t-o /path/to/output.csv (optional; default is stdout)\n");
    fprintf(stderr, "\t-v verbose output (to stderr)\n");
    fprintf(stderr, "\t-h this help\n");
}

int main(int argc, char **argv)
{
    int ch;
    char *input_filename = NULL;
    char *output_filename = NULL;
    
    while ((ch = getopt(argc, argv, "i:o:k:vh")) != -1) {
        switch (ch) {
        case 'i':
            input_filename = optarg;
            break;
        case 'o':
            output_filename = optarg;
            break;
        case 'k':
            parse_output_keys(optarg);
            break;
        case 'v':
            verbose = 1;
            break;
        case 'h':
            usage();
            exit(0);
            break;
        default:
            usage();
            exit(1);
            break;
        }
    }
    
    if (!num_output_keys){
        fprintf(stderr, "no input keys specified\n");
        usage();
        exit(1);
    }
    
    run(input_filename, output_filename);
    return 0;
}
