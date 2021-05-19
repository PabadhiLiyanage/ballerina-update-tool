/*
 * Copyright (c) 2021, WSO2 Inc. (http://wso2.com) All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.ballerinalang.command;

import org.ballerinalang.command.cmd.PullCommand;
import org.ballerinalang.command.exceptions.CommandException;
import org.testng.Assert;
import org.testng.annotations.Test;
import picocli.CommandLine;

/**
 * Test cases for dist pull command.
 *
 * @since 2.0.0
 */
public class PullCommandTest extends CommandTest {

    @Test(expectedExceptions = { CommandException.class })
    public void pullCommandwithoutArgsTest() throws CommandException {
        PullCommand pullCommand = new PullCommand(testStream);
        new CommandLine(pullCommand).parse();
        pullCommand.execute();
    }

    @Test(expectedExceptions = { CommandException.class })
    public void pullCommandwithMultipleArgsTest() throws CommandException {
        PullCommand pullCommand = new PullCommand(testStream);
        new CommandLine(pullCommand).parse("arg1", "arg2");
        pullCommand.execute();
    }

    @Test(expectedExceptions = { CommandException.class })
    public void pullNonExistingDist() throws CommandException {
        PullCommand pullCommand = new PullCommand(testStream);
        new CommandLine(pullCommand).parse("slp9");
        pullCommand.execute();
    }

    @Test
    public void pullCommandHelpTest() {
        PullCommand pullCommand = new PullCommand(testStream);
        new CommandLine(pullCommand).parse("-h");
        pullCommand.execute();
        Assert.assertTrue(outContent.toString().contains("Fetch a given distribution and set it as the active version"));
        writeOutput(outContent.toString());
    }

    @Test
    public void pullDistributionTest() {
        PullCommand pullCommand = new PullCommand(testStream);
        new CommandLine(pullCommand).parse("slp1");
        pullCommand.execute();
        Assert.assertTrue(outContent.toString().contains("Checking whether any latest update tool version is available"));
        Assert.assertTrue(outContent.toString().contains("Fetching the 'slp1' distribution from the remote server"));
        Assert.assertTrue(outContent.toString().contains("Fetching the dependencies for 'slp1' from the remote server"));
        Assert.assertTrue(outContent.toString().contains("successfully set as the active distribution"));

        PullCommand pullCmd = new PullCommand(testStream);
        new CommandLine(pullCmd).parse("latest");
        pullCmd.execute();
        Assert.assertTrue(outContent.toString().contains("Fetching the latest distribution from the remote server"));
        Assert.assertTrue(outContent.toString().contains("successfully set as the active distribution"));

        pullCommand.execute();
        Assert.assertTrue(outContent.toString().contains("is already available locally"));
    }
}
