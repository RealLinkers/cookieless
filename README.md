# cookieless

Super simple way to check for cookieless reflected XSS, but also covers potential reflected path XSS

Script checks if the following characters are encoded on the output:

 - "
 - <
 - '

#### Example use:
echo "https://testurl.com" | cookieless

works great in combination with tomnomnom tools like httprobe
https://github.com/tomnomnom/httprobe

echo "testurl.com" | httprobe | cookieless
