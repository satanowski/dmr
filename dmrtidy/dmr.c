#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>


char** split(char *line) {
    char **result = 0;
    size_t count = 0;
    char *tmp = line;
    char *last_coma = 0;
    char delim[2];
    delim[0] = ';';
    delim[1] = 0;
    
    while (*tmp) {
        if (delim[0] == *tmp) {
            count++;
            last_coma = tmp;
        }
        tmp++;
    }

    count += last_coma < (line + strlen(line) -1);
    count ++;

    result = malloc(sizeof(char*) * count);
    if (result) {
        size_t idx = 0;
        char *token = strtok(line, delim);
        while(token) {
            assert(idx < count);
            *(result + idx++) = strdup(token);
            token = strtok(0, delim);
        }
        *(result + idx) = 0;
    }
    return result;
}

void process_line(char *line) {
    char **tokens;
    tokens = split(line);
    if (tokens) {
        int i;
        for(i=0; *(tokens + i); i++) {
            if (i>0 && i<4) printf("%s;", *(tokens + i));
            if (i==4) { // country fix
                if ( *(tokens+i)[0] == ' ') printf("-Unknown-");
                else if (strcmp(*(tokens+i), "ARE") == 0) printf("United Arab Emirates");
                else if (strcmp(*(tokens+i), "ARG") == 0) printf("Argentina");
                else if (strcmp(*(tokens+i), "AUS") == 0) printf("Australia");
                else if (strcmp(*(tokens+i), "BRA") == 0) printf("Brazil");
                else if (strcmp(*(tokens+i), "BRN") == 0) printf("Brunei Darussalam");
                else if (strcmp(*(tokens+i), "CAN") == 0) printf("Canada");
                else if (strcmp(*(tokens+i), "CHL") == 0) printf("Chile");
                else if (strcmp(*(tokens+i), "CHN") == 0) printf("China");
                else if (strcmp(*(tokens+i), "DNK") == 0) printf("Denmark");
                else if (strcmp(*(tokens+i), "HKG") == 0) printf("Hong Kong");
                else if (strcmp(*(tokens+i), "IDN") == 0) printf("Indonesia");
                else if (strcmp(*(tokens+i), "IND") == 0) printf("India");
                else if (strcmp(*(tokens+i), "ISR") == 0) printf("Israel");
                else if (strcmp(*(tokens+i), "JPN") == 0) printf("Japan");
                else if (strcmp(*(tokens+i), "KOR") == 0) printf("Korea");
                else if (strcmp(*(tokens+i), "KWT") == 0) printf("Kuwait");
                else if (strcmp(*(tokens+i), "MEX") == 0) printf("Mexico");
                else if (strcmp(*(tokens+i), "MYS") == 0) printf("Malaysia");
                else if (strcmp(*(tokens+i), "NLD") == 0) printf("Netherlands");
                else if (strcmp(*(tokens+i), "NZL") == 0) printf("New Zealand");
                else if (strcmp(*(tokens+i), "PAN") == 0) printf("Panama");
                else if (strcmp(*(tokens+i), "PHL") == 0) printf("Philippines");
                else if (strcmp(*(tokens+i), "QAT") == 0) printf("Quatar");
                else if (strcmp(*(tokens+i), "SGP") == 0) printf("Singapore");
                else if (strcmp(*(tokens+i), "THA") == 0) printf("Thailand");
                else if (strcmp(*(tokens+i), "TTO") == 0) printf("Trinidad and Tobago");
                else if (strcmp(*(tokens+i), "TWN") == 0) printf("Taiwan");
                else if (strcmp(*(tokens+i), "URY") == 0) printf("Uruguay");
                else if (strcmp(*(tokens+i), "VEN") == 0) printf("Venezuela");
                else if (strcmp(*(tokens+i), "Argentina Republic") == 0) printf("Argentina");
                else if (strcmp(*(tokens+i), "Belarus/Belorussia") == 0) printf("Belarus");
                else if (strcmp(*(tokens+i), "Bosnia and Hercegovi") == 0) printf("Bosnia and Hercegovina");
                else if (strcmp(*(tokens+i), "Columbia") == 0) printf("Colombia");
                else if (strcmp(*(tokens+i), "Faroe") == 0) printf("Faroe Islands");
                else if (strcmp(*(tokens+i), "Korea S, Republic of") == 0) printf("Korea");
                else if (strcmp(*(tokens+i), "Korea, Republic of") == 0) printf("Korea");
                else if (strcmp(*(tokens+i), "Lithunia") == 0) printf("Lithuania");
                else if (strcmp(*(tokens+i), "Macao, China") == 0) printf("Macao");
                else if (strcmp(*(tokens+i), "Moldava") == 0) printf("Moldova");
                else if (strcmp(*(tokens+i), "St. Vincent & Gren.") == 0) printf("Saint Vincent and the Grenadines");
                else if (strcmp(*(tokens+i), "Swasiland") == 0) printf("Swaziland");
                else if (strcmp(*(tokens+i), "Virgin Islands") == 0) printf("British Virgin Islands");
                else if (strcmp(*(tokens+i), "Virgin Islands, U.S.") == 0) printf("US Virgin Islands");
                else if (strcmp(*(tokens+i), "country") == 0) printf("-Unknown-");
                else printf("%s", *(tokens + i));
            }
            free(*(tokens + i));
        }
        printf("\n");
        free(tokens);
    }
}

int main(int argc, char *argv[]) {
    FILE *fp = stdin;
    char *line = NULL;
    size_t len = 0;
    ssize_t read;

    if (argc == 2 && (strcmp(argv[1], "-h")==0 || strcmp(argv[1], "-H")==0)) {
        printf("\nUsage:\n\n");
        printf("Clean the input stream from incorrect and redundant data:\n");
        printf("  cat contacts.csv | %s > filtered.csv\n", argv[0]);
        return 0;
    }

    if (fp == NULL) exit(1);
    while((read = getline(&line, &len, fp)) != -1) {
        process_line(line);
    }

    if (ferror(fp)) {
        printf("Error: :\n");
    }

    free(line);
    fclose(fp);

    return 0;
}
