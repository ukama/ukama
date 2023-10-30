BEGIN { start = 0; incomment = 0;}
{
    if (start == 0)
    {
        if (incomment == 0)
        {
            if (/^\/\//)
            {
            }
            else if (/^\/\*[^\/]*\*\/$/)
            {
            }
            else if (/^\/\*/)
            {
                incomment = 1
            }
            else
            {
                print $0;
                start = 1;
            }
        }
        else
        {
            if (/\*\//)
            {
                incomment = 0;
            }
        }
    }
    else
    {
        print $0
    }
}
