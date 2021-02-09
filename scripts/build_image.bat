@ECHO OFF
SET n_args=0
for %%x in (%*) do Set /A n_args+=1

IF NOT %n_args%==2 (
    ECHO Usage: build_image ^<dockerfile_dir^> ^<image_name^>
    GOTO:EOF
)
docker build -t %2 %1


