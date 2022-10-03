package ServerBox

import (
	"github.com/YuranIgnatenko/Json"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
)

// enable loop refresh page
func EnableReloadPage(w http.ResponseWriter, sec int) {
	F(w, "<head><meta http-equiv=\"refresh\" content=%v></head>", sec)
}

func AddImage(w http.ResponseWriter, namefile string, width, height int) {
	F(w, "<div><img src=\"%v\" alt=\"image file: %v\" width=\"%v\" height=\"%v\"></div>", namefile, namefile, width, height)
}

func AddCSS(w http.ResponseWriter, namefileCSS string) error {
	data, err := os.ReadFile(namefileCSS)
	if err != nil {
		return err
	}
	code := `<style type="text/css">` + string(data) + `</style>`
	F(w, code)
	return nil
}

// add header in page; (w, "Hello", 3) -> <h3>Hello</h3>
func AddLine(w http.ResponseWriter, line string, size int) error {
	if 0 < size && size < 7 {
		F(w, "<div><h%v>%v</h%v></div>", size, line, size)
		return nil
	}
	lineerror := S("<h%v>%v</h%v>", size, line, size)
	return errors.New(S("error line compile: %v\n", lineerror))
}

// run javascript string
func RunScriptJS(w http.ResponseWriter, codejs, id string) {
	F(w, "<script id=%v type=\"text/javascript\">%v</script>", id, codejs)
}

// run javavscript file
func RunFileJS(w http.ResponseWriter, filename string) error {
	data, er := os.ReadFile(filename)
	if er != nil {
		return errors.New(S("error get file: %v\n", filename))
	}
	codejs := string(data)
	F(w, "<script type=\"text/javascript\">%v</script>", codejs)
	return nil
}

func AddHtmlPage(w http.ResponseWriter, filename string) {
	var temp = template.Must(template.ParseFiles(filename))
	temp.Execute(w, nil)
}

func EnableSupportAjax(w http.ResponseWriter) {
	data := `<script type="text/javascript" src="http://code.jquery.com/jquery-latest.min.js"></script>`
	F(w, data)
}

func ActivateAddressAjax(w http.ResponseWriter, receiver_path, sendParams string) {
	data := `<script>
            $.ajax({
                url:"` + S("%v", receiver_path) + `",
                method:"POST",
                data: {
                    send:"` + S("%v", sendParams) + `",
                },
            });
	</script>`
	F(w, data)
}

func GetAjax(w http.ResponseWriter, r *http.Request, method string) (bool, string) {
	if r.Method == method {
		var err error
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return false, ""
		}
		return true, string(b)
	}
	return false, ""
}

func RunAnimateImageFromPack(
	w http.ResponseWriter,
	canvasId,
	namefile,
	width_frame,
	height_frame,
	width_all,
	height_all,
	count_frames,
	tick string,
) {
	F(w, "<canvas id=\"%v\"></canvas>", canvasId)
	RunScriptJS(w, `
	class Sprite {
    constructor(options) {
        this.ctx = options.ctx;
        this.image = options.image;
        this.frameIndex = 0;
        this.tickCount = 0;
        this.ticksPerFrame = options.ticksPerFrame || 0;
        this.numberOfFrames = options.numberOfFrames || 1;
        this.width = options.width;
        this.height = options.height;
        this.start();
    }
    update() {
        this.tickCount++;
        if (this.tickCount > this.ticksPerFrame) {
            this.tickCount = 0;
            if (this.frameIndex < this.numberOfFrames - 1) {
                this.frameIndex++;
            } else {
                this.frameIndex = 0;
            }
        }
    }
    render() {
        this.ctx.clearRect(0, 0, this.width / this.numberOfFrames, this.height);
        this.ctx.drawImage(
            this.image,
            this.frameIndex * this.width / this.numberOfFrames,
            0,
            this.width / this.numberOfFrames,
            this.height,
            0,
            0,
            this.width / this.numberOfFrames,
            this.height
        )
    }
    start() {
        let loop = () => {
            this.update();
            this.render();
            window.requestAnimationFrame(loop);
        }
        window.requestAnimationFrame(loop);
    }
}
let canvas = document.getElementById('`+canvasId+`');
canvas.width = `+width_frame+`;
canvas.height = `+height_frame+`;
let coinImage = new Image();
coinImage.src = '`+namefile+`';
let sprite = new Sprite({
  ctx: canvas.getContext('2d'),
  image: coinImage,
  width: `+width_all+`,
  height: `+height_all+`,
  numberOfFrames: `+count_frames+`,
  ticksPerFrame: `+tick+`,
})
	`,
		"animateId")
}

func AddCanvasImage(
	w http.ResponseWriter,
	canvasId,
	imageId,
	namefile,
	width_frame,
	height_frame string,
) {
	RunScriptJS(w, `
	
let canvas = document.getElementById('`+canvasId+`');
canvas.width = `+width_frame+`;
canvas.height = `+height_frame+`;
let coinImage = new Image();
coinImage.src = '`+namefile+`';
canvas.drawImage(coinImage,0,0)
})
	`, imageId)
}

func ReturnRequestJSON(w http.ResponseWriter, DataStruct any) error {
	Data, err := Json.StructToJson(DataStruct)
	if err != nil {
		return err
	}
	jData, err := json.Marshal(Data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(string(jData))
	w.Write(jData)
	fmt.Println(w.Header().Get("Content-Type"))

	return nil
}
