/*
Sources for the audio dsp fw isn't included, despite the kernel module claiming to have a Dual BSD/GPL license:
# modinfo /lib/modules/ismdaudio.ko
filename:       /lib/modules/ismdaudio.ko
license:        Dual BSD/GPL
..

So we need to pull the data from the original kernel module.
*/

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

/*
Constants come from:

./build/target/bin/i686-cm-linux-objdump -t lib/modules/ismdaudio.ko  | grep kernel
00000000 l    df *ABS*	00000000 audio_dsp_fw_kernel_dsp0.c
00000000 l    df *ABS*	00000000 audio_dsp_fw_kernel_dsp1.c
000051c0 g     O .rodata	00200000 audio_fw_kernel_dsp0
002051c0 g     O .rodata	00200000 audio_fw_kernel_dsp1
./build/target/bin/i686-cm-linux-objdump -h lib/modules/ismdaudio.ko  | grep .rodata
  1 .rodata       00407347  00000000  00000000  00064180  2**5

*/
const (
	dsp_size      = 0x200000
	dsp0_off      = 0x51c0
	dsp1_off      = 0x2051c0
	rodata_offset = 0x64180
	input_ko      = "./lib/modules/ismdaudio.ko"
	output_c0     = "./build/target/src/audio/core/audio_dsp_fw_kernel_dsp0.c"
	output_c1     = "./build/target/src/audio/core/audio_dsp_fw_kernel_dsp1.c"
)

type bin2C struct {
	in      bytes.Buffer
	first   bool
	written int64
}

func NewBin2C() io.ReadWriter {
	return &bin2C{}
}

func (bc *bin2C) Write(data []byte) (n int, err error) {
	l := len(data)
	if !bc.first {
		bc.first = true
		bc.in.WriteString(fmt.Sprintf("%#02x", data[0]))
		data = data[1:]
	}
	for i, b := range data {
		bc.written++
		bc.in.WriteRune(',')
		if bc.written%16 == 0 {
			bc.in.WriteRune('\n')
		}
		if _, err := bc.in.WriteString(fmt.Sprintf("%#02x", b)); err != nil {
			return i, err
		}
	}
	return l, nil
}

func (b *bin2C) Read(data []byte) (n int, err error) {
	return b.in.Read(data)
}

func conv(f io.Reader, outfn, symbol string) {
	os.MkdirAll(path.Dir(outfn), 0755)
	out1, err := os.Create(outfn)
	if err != nil {
		log.Fatalln(err)
	}
	defer out1.Close()
	out1.WriteString(fmt.Sprintf("const char %s[] = {\n", symbol))
	f1 := io.LimitReader(f, dsp_size)
	bc := NewBin2C()
	if _, err := io.Copy(bc, f1); err != nil {
		log.Fatalln(err)
	}
	if _, err := io.Copy(out1, bc); err != nil {
		log.Fatalln(err)
	}
	out1.WriteString("};\n")
}
func main() {
	f, err := os.Open(input_ko)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	_, err = f.Seek(rodata_offset+dsp0_off, 0)
	if err != nil {
		log.Fatalln(err)
	}
	conv(f, output_c0, "audio_dsp_fw_kernel_dsp0")
	_, err = f.Seek(rodata_offset+dsp1_off, 0)
	if err != nil {
		log.Fatalln(err)
	}
	conv(f, output_c1, "audio_dsp_fw_kernel_dsp1")
}
