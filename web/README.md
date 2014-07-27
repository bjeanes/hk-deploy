# hk-deploy web

Eventually, a small HTTP that serves out short-lived policies for uploading a
single file to an S3 bucket and a short-lived URL for accessing that file.

It will be used by the `hk` plugin at the root of this repository (`hk-deploy`)
to upload a directory of source code as a `.tgz` file so the Heroku Build API
can [create a build using the source](https://devcenter.heroku.com/articles/build-and-release-using-the-api#creating-builds).

The bucket that stores this code has the minimum allowable life cycle (24 hours)
for automatically deleting files. It only accepts `.tgz` files that are less than
200 MB\*. The bucket and all items are set to private so the URL provided by this
service is the only way to retrieve the file and it will only work for 30 minutes.

\* restriction not currently implemented

## Testing

``` sh-session
$ hk env -a hk-deploy | xargs -J __ env __ go run *.go curl
curl https://hk-deploy.s3.amazonaws.com/ \
  -Fkey="20140727/9c5ba3ed-91fe-4ab3-9c4a-705b3dab9bb4.tgz" \
  -Facl=private -FContent-Type=application/x-compressed \
  -FX-Amz-Credential="AKIAI73CN3SXNGBBBYEA/20140727/us-east-1/s3/aws4_request" \
  -FX-Amz-Algorithm=AWS4-HMAC-SHA256 \
  -FX-Amz-Date="20140727T020425Z" \
  -FPolicy="ewoJImV4cGlyYXRpb24iOiAiMjAxNC0wNy0yN1QwMjoxNDoyNVoiLAoJImNvbmRpdGlvbnMiOiBbCgkJeyJhY2wiOiAicHJpdmF0ZSIgfSwKCQl7ImJ1Y2tldCI6ICJoay1kZXBsb3kiIH0sCgkJeyJrZXkiOiAiMjAxNDA3MjcvOWM1YmEzZWQtOTFmZS00YWIzLTljNGEtNzA1YjNkYWI5YmI0LnRneiJ9LAoJCXsieC1hbXotZGF0ZSI6ICIyMDE0MDcyN1QwMjA0MjVaIn0sCgkJeyJ4LWFtei1jcmVkZW50aWFsIjogIkFLSUFJNzNDTjNTWE5HQkJCWUVBLzIwMTQwNzI3L3VzLWVhc3QtMS9zMy9hd3M0X3JlcXVlc3QifSwKCQl7IngtYW16LWFsZ29yaXRobSI6ICJBV1M0LUhNQUMtU0hBMjU2In0sCgkJeyJjb250ZW50LXR5cGUiOiAiYXBwbGljYXRpb24veC1jb21wcmVzc2VkIn0KCV0KfQ==" \
  -FX-Amz-Signature="fc32665b200bfbbf02cc8a46ae3d32ff9923242ee2d549f79eb14b6db15db43e" \
  -Ffile=@some-tarred-code.tgz
```
