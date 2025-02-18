const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Virtual Classroom</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 20px;
            text-align: center;
        }
        .container {
            max-width: 600px;
            background: #ffffff;
            padding: 25px;
            border-radius: 12px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            margin: auto;
        }
        .header {
            font-size: 24px;
            font-weight: bold;
            color: #0056b3;
        }
        .content {
            font-size: 16px;
            color: #333;
            margin-top: 15px;
            line-height: 1.6;
        }
        .footer {
            font-size: 14px;
            color: #666;
            margin-top: 20px;
            border-top: 1px solid #ddd;
            padding-top: 10px;
        }
        .cta-button {
            background: #007bff;
            color: #fff;
            padding: 12px 20px;
            text-decoration: none;
            font-size: 16px;
            border-radius: 6px;
            display: inline-block;
            margin-top: 15px;
        }
        .cta-button:hover {
            background: #0056b3;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2 class="header">Welcome to Virtual Classroom, {{ .Name }}! 👋</h2>
        <p class="content">
            We're thrilled to have you as part of our learning community! 🚀<br>
            Our platform provides interactive learning experiences, personalized resources, and a collaborative space for growth.
        </p>
        <a href="https://your-virtual-classroom.com/dashboard" class="cta-button">Go to Dashboard</a>
        <p class="footer">
            Best Regards, <br>
            <strong>Virtual Classroom Team</strong> <br>
            <a href="https://your-virtual-classroom.com">your-virtual-classroom.com</a>
        </p>
    </div>
</body>
</html>
`
