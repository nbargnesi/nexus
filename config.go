// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// See http://formwork-io.github.io/ for more.

package main

import (
	"errors"
	"fmt"
	toml "github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type rail struct {
	Name    string
	Pattern string
	Ingress int
	Egress  int
}

type rails struct {
	Rail []rail
}

func ReadConfigFile(path string) ([]rail, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("configuration file error: %s", err.Error()))
	}

	var rails rails
	_, err = toml.Decode(string(data), &rails)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("configuration file error: %s", err.Error()))
	}

	_, err = validateRails(rails.Rail)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("configuration file error: " + err.Error()))
	}

	return rails.Rail, nil
}

func ReadEnvironment() ([]rail, error) {
	GL_RAIL_NAME_TMPL := "GL_RAIL_%d_NAME"
	GL_RAIL_PATTERN_TMPL := "GL_RAIL_%d_PATTERN"
	GL_RAIL_INGRESS_TMPL := "GL_RAIL_%d_INGRESS_PORT"
	GL_RAIL_EGRESS_TMPL := "GL_RAIL_%d_EGRESS_PORT"

	var rails []rail
	index := 0
	for {
		name, err := getenv(fmt.Sprintf(GL_RAIL_NAME_TMPL, index))
		if err != nil {
			if index != 0 {
				break
			}
			return nil, err
		}
		if name == "" {
			break
		}

		pattern, err := getenv(fmt.Sprintf(GL_RAIL_PATTERN_TMPL, index))
		if err != nil {
			return nil, err
		}
		ingress, err := getenv(fmt.Sprintf(GL_RAIL_INGRESS_TMPL, index))
		if err != nil {
			return nil, err
		}
		egress, err := getenv(fmt.Sprintf(GL_RAIL_EGRESS_TMPL, index))
		if err != nil {
			return nil, err
		}
		ingress_port, err := asPort(ingress)
		if err != nil {
			return nil, err
		}
		egress_port, err := asPort(egress)
		if err != nil {
			return nil, err
		}

		rails = append(rails, rail{
			Name:    name,
			Pattern: pattern,
			Ingress: ingress_port,
			Egress:  egress_port,
		})

		index++
	}

	_, err := validateRails(rails)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("configuration file error: " + err.Error()))
	}

	return rails, nil
}

func validateRails(rails []rail) (*rail, error) {
	re := regexp.MustCompile("[[:punct:]]|[[:space:]]")
	for i, rail := range rails {
		normalizedPattern := strings.ToLower(re.ReplaceAllLiteralString(rail.Pattern, ""))
		if normalizedPattern != "pubsub" && normalizedPattern != "reqrep" {
			msg := fmt.Sprintf("invalid pattern \"%s\" for rail \"%s\"", rail.Pattern, rail.Name)
			return &rail, errors.New(msg)
		}
		rails[i].Pattern = normalizedPattern
	}

	return nil, nil
}

func getenv(env string) (string, error) {
	_env := os.Getenv(env)
	if len(_env) == 0 {
		return "", errors.New(fmt.Sprintf("no %s is set", env))
	}
	return _env, nil
}

func asPort(env string) (int, error) {
	port, err := strconv.Atoi(env)
	if err != nil {
		die("invalid port: %s", env)
		return -1, errors.New(fmt.Sprintf("invalid port: %v - %s", env, err.Error()))
	} else if port < 1 || port > 65535 {
		die("invalid port: %s", env)
		return -1, errors.New(fmt.Sprintf("invalid port: %v - %s", env, err.Error()))
	}
	return port, nil
}
