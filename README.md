# cropper

go implimentation of cropping images

## Usage

```
C:\work> cropper.exe
Usage: cropper.exe [ -conf conf.json ] infile_pattern [ out_dir ]
```

you can treat wild card in infile_pattern.

```
C:\work> cropper.exe /tmp/*.png out
```

If you ommit -conf option, default conf.json file is refered. 

Default conf.json file must be located at the same path of cropper.exe.

## Config

[conf.json](conf.json)

```
{
    "p1":[29, 310],
    "p2":[709, 820],
    "out_dir": "out"
}
```

cropper crops area of between p1 and p2 from input image file and saves in out_dir.

+ "p1" is the coodinate of top-left point.
+ "p2" is the coodinate of bottom-right point.
+ "out_dir" is default out directory.
