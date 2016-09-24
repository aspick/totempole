package totempole

import (
    "io/ioutil"
    "runtime"
    "path"

    "gopkg.in/yaml.v2"
)

type Pole struct {
    Name    string
    Ps      string
    Cmd     string
    Sh      string
    Pwd     string
    Workers int
}

type Totemfile struct {
    Meta string
    Daemons []Pole
}

func readTotemfile() (Totemfile, error) {
    // get default Totemfile path
    _, filename, _, _ := runtime.Caller(1)
    path := path.Join(path.Dir(path.Dir(filename)), "Totemfile.yml")

    buf, e := ioutil.ReadFile(path)
    totemfile := Totemfile{}

    if e != nil {
        return totemfile, e
    }

    e = yaml.Unmarshal(buf, &totemfile)
    if e != nil {
        return totemfile, e
    }

    return totemfile, nil
}
