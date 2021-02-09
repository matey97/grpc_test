@ECHO OFF
SET n_args=0
for %%x in (%*) do Set /A n_args+=1

IF NOT %n_args%==3 (
    ECHO Usage: deploy_image.bat ^<local_port^> ^<docker_port^> ^<image_name^>
    GOTO:EOF
)

docker run -dp %1:%2 ^
 -e PORT=%2 ^
 -e GOOGLE_APPLICATION_CREDENTIALS=/tmp/keys/google_services.json %3